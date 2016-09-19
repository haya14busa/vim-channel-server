"=============================================================================
" FILE: plugin/channel_server.vim
" AUTHOR: haya14busa
" License: MIT license
"=============================================================================
scriptencoding utf-8
if expand('%:p') ==# expand('<sfile>:p')
  unlet! g:loaded_channel_server
endif
if exists('g:loaded_channel_server')
  finish
endif
let g:loaded_channel_server = 1
let s:save_cpo = &cpo
set cpo&vim

if get(g:, 'vim_channel_server#enable_at_startup', v:true)
  call channel_server#serve(get(g:, 'vim_channel_server#addr', ''))
endif

let &cpo = s:save_cpo
unlet s:save_cpo
" __END__
" vim: expandtab softtabstop=2 shiftwidth=2 foldmethod=marker
