"=============================================================================
" FILE: autoload/channel_server.vim
" AUTHOR: haya14busa
" License: MIT license
"=============================================================================
scriptencoding utf-8
let s:save_cpo = &cpo
set cpo&vim

let g:vim_channel_server#debug = v:false
let g:vim_channel_server#errors = []

call ch_logfile('/tmp/channellog', 'w')

" Get the directory separator.
function! s:separator() abort
  return fnamemodify('.', ':p')[-1 :]
endfunction

let s:base = expand('<sfile>:p:h:h')
let s:cmd = s:base . s:separator() . fnamemodify(s:base, ':t')

function! s:err_cb(ch, msg) abort
  let g:vim_channel_server#errors += [a:msg]
  if g:vim_channel_server#debug
    echom string(a:ch) a:msg
  endif
endfunction

let s:option = {
\   'in_mode': 'json',
\   'out_mode': 'json',
\
\   'err_cb': function('s:err_cb'),
\ }

function! channel_server#serve(addr) abort
  if exists('g:vim_channel_server#job') && job_status(g:vim_channel_server#job) ==# 'run'
    throw 'vim-channel-server: server is running already'
  endif
  if exepath(s:cmd) ==# ''
    throw 'vim-channel-server: executable not found. Please run $ make'
  endif

  unlockvar! g:vim_channel_server#cmd
  let g:vim_channel_server#cmd = s:cmd
  lockvar! g:vim_channel_server#cmd

  let cmd = [g:vim_channel_server#cmd]
  if a:addr !=# ''
    let cmd += ['-addr', a:addr]
  endif

  let g:vim_channel_server#job = job_start(cmd, s:option)
endfunction

let &cpo = s:save_cpo
unlet s:save_cpo
" __END__
" vim: expandtab softtabstop=2 shiftwidth=2 foldmethod=marker
