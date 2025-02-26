import { useImperativeHandle, useState } from 'react';
import type { RefObject } from 'react';

export interface DialogRef {
  openDialog: () => void;
  closeDialog: () => void;
}

export function useDialog(ref: RefObject<DialogRef | null>, onClose?: () => void) {
  const [isOpen, setIsOpen] = useState<boolean>(false);

  const handleClose = () => {
    setIsOpen(false);
    onClose?.();
  };

  useImperativeHandle(ref, () => ({
    openDialog: () => setIsOpen(true),
    closeDialog: handleClose,
  }));

  return {
    isOpen,
    handleClose,
  };
}
