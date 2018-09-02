set nocompatible
filetype off
" ================ General Config ====================

set number                      "Line numbers are good
set backspace=indent,eol,start  "Allow backspace in insert mode
set history=1000                "Store lots of :cmdline history
set showcmd                     "Show incomplete cmds down the bottom
set showmode                    "Show current mode down the bottom
set gcr=a:blinkon0              "Disable cursor blink
set visualbell                  "No sounds
set autoread                    "Reload files changed outside vim

set hidden

syntax on

set noswapfile
set nobackup
set nowb

set smarttab
set shiftwidth=8
set softtabstop=8
set tabstop=8
set noexpandtab

set wrap
set linebreak

" Enable completion where available.
let g:ale_completion_enabled = 1
let g:ale_echo_msg_format = '[%linter%] %s'
let g:ale_javascript_eslint_options = {"indent": "off", "object-curly-spacing" : "off", "space-in-parens" : "off"}


set rtp+=~/.vim/bundle/Vundle.vim
call vundle#begin()
Plugin 'VundleVim/Vundle.vim'
Plugin 'octol/vim-cpp-enhanced-highlight'
Plugin 'micha/vim-colors-solarized'
Plugin 'Chiel92/vim-autoformat'
Plugin 'alvan/vim-closetag'
Plugin 'w0rp/ale'
call vundle#end()
execute pathogen#infect()
filetype plugin indent on
if has('gui_running')
else
let g:solarized_termcolors=256
endif
syntax enable
set background=dark
colorscheme solarized
