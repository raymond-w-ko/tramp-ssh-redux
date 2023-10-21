;;; tsr-repl.el --- -*- lexical-binding: t -*-
;;; Commentary:
;;; This file is for development and quick testing of the go binaries in this project

;;; Code:

(require 'tramp)

(defun tsr/get-parent-directory ()
  ""
  (file-name-directory (or load-file-name buffer-file-name)))

(defvar tsr/client-path
  (expand-file-name "bin/tramp-ssh-redux-client" (tsr/get-parent-directory)))

(defun tsr/create-base-cmd (tramp-buf-name)
  "TODO"
  (let* ((tokens (tramp-dissect-file-name tramp-buf-name))
         (method (tramp-file-name-method tokens))
         (user (tramp-file-name-user tokens))
         (host (tramp-file-name-host tokens))
         (path (tramp-file-name-localname tokens)))
    `(:method ,method :user ,user :host ,host :path ,path)))

(defun tsr/test-client-1 ()
  "TODO"
  (let* ((ssh-path "/ssh:rko@localhost:~/src/tramp-ssh-redux/Makefile")
         (base-cmd (tsr/create-base-cmd ssh-path))
         (cmd (append base-cmd '(:command "echo" :oneshot t)))
         (out-buf (generate-new-buffer (generate-new-buffer-name " *tsr/client-out*"))))
    (call-process tsr/client-path nil out-buf nil (json-serialize cmd))

    (buffer-disable-undo out-buf)
    (with-current-buffer out-buf
      (message (buffer-string)))
    (kill-buffer out-buf)))

;; (tsr/test-client-1)

(provide 'tsr-repl)
;;; tsr-repl.el ends here
