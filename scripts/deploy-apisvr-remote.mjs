/**
 * 部署 coresvr（内嵌 syssvr direct 模式）到远程服务器
 *
 * 用法:
 *   $env:DEPLOY_PASSWORD="***"
 *   $env:DEPLOY_UNIPUSH_SECRET="可选，写入远程 sys.yaml"
 *   $env:DEPLOY_UPLOAD_SYSYAML="1"  # 默认不上传 sys.yaml，避免覆盖远程 UniPush
 *   node scripts/deploy-apisvr-remote.mjs
 */
import { spawnSync } from 'node:child_process'
import { existsSync, mkdirSync, cpSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'
import { Client } from 'ssh2'

const __dirname = dirname(fileURLToPath(import.meta.url))
const root = join(__dirname, '..')

const host = process.env.DEPLOY_HOST || '47.94.112.109'
const user = process.env.DEPLOY_USER || 'root'
const password = process.env.DEPLOY_PASSWORD || ''
const remotePath = process.env.DEPLOY_PATH || '/root/run/core'
const unipushSecret = process.env.DEPLOY_UNIPUSH_SECRET || ''

if (!password) {
  console.error('请设置环境变量 DEPLOY_PASSWORD')
  process.exit(1)
}

function run(cmd, args, env = {}) {
  const r = spawnSync(cmd, args, {
    cwd: root,
    stdio: 'inherit',
    env: { ...process.env, ...env },
    shell: process.platform === 'win32',
  })
  if (r.status !== 0) {
    process.exit(r.status ?? 1)
  }
}

function sshExec(conn, cmd) {
  return new Promise((resolve, reject) => {
    conn.exec(cmd, (err, stream) => {
      if (err) return reject(err)
      let out = ''
      const timeout = setTimeout(() => {
        reject(new Error(`remote command timeout: ${cmd}`))
      }, 120000)
      let finished = false
      const done = (ok, val) => {
        if (finished) return
        finished = true
        clearTimeout(timeout)
        if (ok) resolve(val)
        else reject(val)
      }
      stream
        .on('close', (code) => {
          if (code === 0 || code === null) done(true, out)
          else done(false, new Error(`remote exit ${code}: ${out}`))
        })
        .on('data', (d) => {
          process.stdout.write(d)
          out += d.toString()
        })
      stream.on('end', () => {
        done(true, out)
      })
      stream.stderr.on('data', (d) => process.stderr.write(d))
    })
  })
}

function sftpPut(sftp, local, remote) {
  return new Promise((resolve, reject) => {
    sftp.fastPut(local, remote, (err) => (err ? reject(err) : resolve()))
  })
}

async function main() {
  const cmdDir = join(root, 'cmd')
  const binLocal = join(cmdDir, 'coresvr')
  const etcLocal = join(cmdDir, 'etc')
  const prebuilt = process.env.CORESVR_PREBUILT || ''

  if (prebuilt && existsSync(prebuilt)) {
    console.log('>> 使用预编译二进制:', prebuilt)
    mkdirSync(cmdDir, { recursive: true })
    cpSync(prebuilt, binLocal)
  } else {
    console.log('>> 编译 coresvr linux/amd64 ...')
    run('go', [
      'build',
      '-tags',
      'no_k8s',
      '-o',
      './cmd/coresvr',
      './service/apisvr',
    ], {
      GOOS: 'linux',
      GOARCH: 'amd64',
      CGO_ENABLED: '0',
      GOTMPDIR: process.env.GOTMPDIR || join(root, '..', '..', 'web-260530', 'tmp-go'),
    })
  }

  if (!existsSync(binLocal)) {
    console.error('编译产物不存在:', binLocal)
    process.exit(1)
  }

  mkdirSync(etcLocal, { recursive: true })
  cpSync(join(root, 'service/apisvr/etc'), etcLocal, { recursive: true, force: true })

  const conn = new Client()
  await new Promise((resolve, reject) => {
    conn
      .on('ready', resolve)
      .on('error', reject)
      .connect({ host, port: 22, username: user, password, readyTimeout: 20000 })
  })

  console.log('>> 上传 coresvr 与配置 ...')
  await sshExec(conn, `id; pwd; ls -ld /root /root/run /root/run/core 2>/dev/null || true; df -h /root || true`)
  await sshExec(conn, `mkdir -p '${remotePath}' '${remotePath}/etc'`)
  const sftp = await new Promise((resolve, reject) => {
    conn.sftp((err, s) => (err ? reject(err) : resolve(s)))
  })
  console.log('>> 上传新二进制到 /tmp ...')
  await sftpPut(sftp, binLocal, `/tmp/coresvr.new`)

  const uploadSysYaml = process.env.DEPLOY_UPLOAD_SYSYAML === '1'
  if (uploadSysYaml) {
    await sftpPut(sftp, join(etcLocal, 'sys.yaml'), `/tmp/sys.yaml.new`)
    await sshExec(conn, `mv -f /tmp/coresvr.new '${remotePath}/coresvr' && chmod +x '${remotePath}/coresvr'`)
    await sshExec(conn, `mv -f /tmp/sys.yaml.new '${remotePath}/etc/sys.yaml'`)
  } else {
    console.log('>> 跳过 sys.yaml 上传（远程 UniPush 配置不会被覆盖）；需上传时设 DEPLOY_UPLOAD_SYSYAML=1')
    await sshExec(conn, `mv -f /tmp/coresvr.new '${remotePath}/coresvr' && chmod +x '${remotePath}/coresvr'`)
  }

  if (unipushSecret) {
    console.log('>> 写入远程 UniPush Secret（sys.yaml）...')
    const esc = unipushSecret.replace(/'/g, `'\\''`)
    await sshExec(
      conn,
      `sed -i "s/^  Secret:.*/  Secret: ${esc}/" '${remotePath}/etc/sys.yaml' && sed -i "s/^  Enabled:.*/  Enabled: true/" '${remotePath}/etc/sys.yaml'`,
    )
  } else {
    console.log('>> 未设置 DEPLOY_UNIPUSH_SECRET：远程 UniPush 保持 sys.yaml 当前值')
  }

  console.log('>> 重启 coresvr（run.sh，与线上一致 ./coresvr core）...')
  await sshExec(
    conn,
    `cd '${remotePath}' && bash run.sh && sleep 5 && pgrep -af '[./]coresvr' | grep -v pgrep | head -3 || (echo 'run.sh failed'; tail -20 coresvr.log 2>/dev/null)`,
  )

  conn.end()
  console.log('>> 部署完成')
}

main().catch((e) => {
  console.error(e)
  process.exit(1)
})
