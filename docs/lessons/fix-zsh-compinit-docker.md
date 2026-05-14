# Fix zsh compinit Docker errors

## 现象

```
source ~/.zshrc
compinit:527: no such file or directory: /opt/homebrew/share/zsh/site-functions/_docker
compinit:527: no such file or directory: /opt/homebrew/share/zsh/site-functions/_docker-compose
compinit:shift:529: shift count must be <= $#
```

## 根因

Homebrew 安装 Docker Desktop 时，在 `/opt/homebrew/share/zsh/site-functions/` 下创建了指向 Docker.app 的符号链接。卸载或移动 Docker.app 后，这些符号链接变成悬空链接（dangling symlinks），compinit 读取时报错。`shift` 错误是前两个错误的连锁反应。

## 排查步骤

1. 先尝试清除 zsh 补全缓存（可能无效）：
   ```bash
   rm -f ~/.zcompdump*
   exec zsh
   ```

2. 若仍报错，检查是否存在悬空符号链接：
   ```bash
   ls -la /opt/homebrew/share/zsh/site-functions/_docker*
   ```

   输出类似：
   ```
   _docker -> /Applications/Docker.app/Contents/Resources/etc/docker.zsh-completion
   _docker-compose -> /Applications/Docker.app/Contents/Resources/etc/docker-compose.zsh-completion
   ```

## 修复

删除悬空符号链接：

```bash
rm /opt/homebrew/share/zsh/site-functions/_docker /opt/homebrew/share/zsh/site-functions/_docker-compose
exec zsh
```

## 通用思路

zsh `compinit` 报 "no such file or directory" 时，排查路径：

1. 清缓存（`rm -f ~/.zcompdump*`）
2. 检查 `site-functions/` 下的悬空符号链接
3. 删除或重建对应的符号链接
