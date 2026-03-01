#!/bin/bash
set -e

apt-get update
apt-get install -y openjdk-21-jdk maven make tmux

# tmux config
cat > /root/.tmux.conf <<'TMUXCONF'
unbind C-b
set-option -g prefix C-a
bind-key C-a send-prefix

# Split pane horizontally with Ctrl-a [
bind-key [ split-window -h -c "#{pane_current_path}"

# Split pane vertically with Ctrl-a ]
bind-key ] split-window -v -c "#{pane_current_path}"

# Close current pane with Ctrl-a =
bind-key = kill-pane

# Ensure new windows also open in the same directory
bind-key c new-window -c "#{pane_current_path}"

bind -n C-k send-keys C-l \;;

set -g mouse on
setw -g mode-keys vi

TMUXCONF

# Auto-launch tmux on login
echo 'if command -v tmux &>/dev/null && [ -z "$TMUX" ]; then tmux new-session -A -s main; fi' >> /root/.bashrc

git clone -b events-service https://github.com/Dan-Sones/prism.git /root/prism

cd /root/prism/services