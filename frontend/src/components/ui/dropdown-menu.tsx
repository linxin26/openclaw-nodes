import * as DropdownMenuPrimitive from "@radix-ui/react-dropdown-menu"

import { cn } from "@/lib/utils"

export const DropdownMenu = DropdownMenuPrimitive.Root
export const DropdownMenuTrigger = DropdownMenuPrimitive.Trigger

export const DropdownMenuContent = ({
  className,
  ...props
}: DropdownMenuPrimitive.DropdownMenuContentProps) => (
  <DropdownMenuPrimitive.Portal>
    <DropdownMenuPrimitive.Content
      className={cn("z-50 min-w-[10rem] rounded-2xl border border-border bg-white p-1 shadow-panel", className)}
      sideOffset={6}
      {...props}
    />
  </DropdownMenuPrimitive.Portal>
)

export const DropdownMenuItem = ({
  className,
  ...props
}: DropdownMenuPrimitive.DropdownMenuItemProps) => (
  <DropdownMenuPrimitive.Item
    className={cn("cursor-pointer rounded-xl px-3 py-2 text-sm outline-none focus:bg-secondary", className)}
    {...props}
  />
)
