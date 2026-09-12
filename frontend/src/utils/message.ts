/**
 * 简单的消息提示工具 (替代 ElMessage)
 * 使用原生 alert/confirm 或自定义 Toast
 */

interface MessageOptions {
  message: string
  type?: 'success' | 'error' | 'warning' | 'info'
  duration?: number
  closable?: boolean
}

// Toast 容器
let toastContainer: HTMLDivElement | null = null
let toastId = 0

function getToastContainer(): HTMLDivElement {
  if (!toastContainer) {
    toastContainer = document.createElement('div')
    toastContainer.id = 'tail-toast-container'
    toastContainer.className = 'fixed top-4 right-4 z-[9999] flex flex-col gap-2'
    document.body.appendChild(toastContainer)
  }
  return toastContainer
}

function showMessage(options: MessageOptions) {
  const { message, type = 'info', duration = 3000, closable = true } = options
  
  const container = getToastContainer()
  const id = ++toastId
  
  const toast = document.createElement('div')
  toast.id = `tail-toast-${id}`
  
  // 颜色样式
  const typeStyles: Record<string, string> = {
    success: 'bg-emerald-50 border-emerald-500 text-emerald-800',
    error: 'bg-red-50 border-red-500 text-red-800',
    warning: 'bg-amber-50 border-amber-500 text-amber-800',
    info: 'bg-blue-50 border-blue-500 text-blue-800',
  }
  
  const iconSvg: Record<string, string> = {
    success: '<svg class="w-5 h-5 text-emerald-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>',
    error: '<svg class="w-5 h-5 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>',
    warning: '<svg class="w-5 h-5 text-amber-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg>',
    info: '<svg class="w-5 h-5 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>',
  }
  
  toast.className = `flex items-center gap-3 px-4 py-3 rounded-lg border shadow-sm ${typeStyles[type] || typeStyles.info}`
  
  let html = `${iconSvg[type] || iconSvg.info}<span class="text-sm font-medium">${message}</span>`
  
  if (closable) {
    html += `<button onclick="document.getElementById('tail-toast-${id}')?.remove()" class="ml-2 text-current opacity-50 hover:opacity-100">&times;</button>`
  }
  
  toast.innerHTML = html
  
  // 添加进入动画
  toast.style.opacity = '0'
  toast.style.transform = 'translateX(100%)'
  toast.style.transition = 'all 0.3s ease'
  
  container.appendChild(toast)
  
  // 触发动画
  requestAnimationFrame(() => {
    toast.style.opacity = '1'
    toast.style.transform = 'translateX(0)'
  })
  
  // 自动消失
  if (duration > 0) {
    setTimeout(() => {
      removeToast(toast)
    }, duration)
  }
  
  return toast
}

function removeToast(toast: HTMLElement) {
  toast.style.opacity = '0'
  toast.style.transform = 'translateX(100%)'
  setTimeout(() => {
    toast.remove()
  }, 300)
}

// Confirm 对话框接口
interface ConfirmOptions {
  title?: string
  message?: string
  type?: 'warning' | 'error' | 'info'
  confirmButtonText?: string
  cancelButtonText?: string
  center?: boolean
}

// Confirm 对话框状态
let confirmState: {
  resolve: (value: boolean) => void
  container: HTMLDivElement
} | null = null

function showConfirm(options: ConfirmOptions): Promise<boolean> {
  return new Promise((resolve) => {
    const {
      title = '确认',
      message: confirmMessage,
      type = 'warning',
      confirmButtonText = '确定',
      cancelButtonText = '取消',
    } = options

    // 创建遮罩层
    const overlay = document.createElement('div')
    overlay.className = 'fixed inset-0 z-[10000] flex items-center justify-center'
    
    // 背景遮罩
    const backdrop = document.createElement('div')
    backdrop.className = 'absolute inset-0 bg-black/50 backdrop-blur-sm'
    backdrop.onclick = () => cleanup()
    
    // 对话框
    const dialog = document.createElement('div')
    dialog.className = 'relative bg-white rounded-xl shadow-2xl w-full max-w-md mx-4 animate-fade-in'
    
    // 根据类型选择图标和颜色
    const typeConfig: Record<string, { color: string; bg: string; icon: string }> = {
      warning: {
        color: 'text-amber-800',
        bg: 'bg-amber-50',
        icon: '<svg class="w-6 h-6 text-amber-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg>'
      },
      error: {
        color: 'text-red-800',
        bg: 'bg-red-50',
        icon: '<svg class="w-6 h-6 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>'
      },
      info: {
        color: 'text-blue-800',
        bg: 'bg-blue-50',
        icon: '<svg class="w-6 h-6 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>'
      }
    }

    const config = typeConfig[type] || typeConfig.warning

    dialog.innerHTML = `
      <div class="p-6">
        <h3 class="text-xl font-semibold text-gray-900 mb-4">${title}</h3>
        <div class="${config.bg} border border-opacity-50 rounded-lg p-4 mb-6">
          <div class="flex items-start gap-3">
            ${config.icon}
            <div class="${config.color} text-sm">${confirmMessage}</div>
          </div>
        </div>
        <div class="flex justify-end gap-3">
          <button class="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors text-sm font-medium" id="confirm-cancel-btn">
            ${cancelButtonText}
          </button>
          <button class="px-4 py-2 bg-${type === 'error' ? 'red' : type === 'warning' ? 'amber' : 'blue'}-600 text-white rounded-lg hover:bg-${type === 'error' ? 'red' : type === 'warning' ? 'amber' : 'blue'}-700 transition-colors text-sm font-medium" id="confirm-ok-btn">
            ${confirmButtonText}
          </button>
        </div>
      </div>
    `

    overlay.appendChild(backdrop)
    overlay.appendChild(dialog)
    document.body.appendChild(overlay)

    // 按钮事件
    dialog.querySelector('#confirm-cancel-btn')?.addEventListener('click', () => {
      cleanup()
      resolve(false)
    })

    dialog.querySelector('#confirm-ok-btn')?.addEventListener('click', () => {
      cleanup()
      resolve(true)
    })

    function cleanup() {
      document.body.removeChild(overlay)
    }

    // ESC 键关闭
    function onKeydown(e: KeyboardEvent) {
      if (e.key === 'Escape') {
        cleanup()
        resolve(false)
        document.removeEventListener('keydown', onKeydown)
      }
    }
    document.addEventListener('keydown', onKeydown)
  })
}

// Message 对象
export const message = {
  success(msg: string, options?: Partial<MessageOptions>) {
    return showMessage({ message: msg, type: 'success', ...options })
  },
  error(msg: string, options?: Partial<MessageOptions>) {
    return showMessage({ message: msg, type: 'error', ...options })
  },
  warning(msg: string, options?: Partial<MessageOptions>) {
    return showMessage({ message: msg, type: 'warning', ...options })
  },
  info(msg: string, options?: Partial<MessageOptions>) {
    return showMessage({ message: msg, type: 'info', ...options })
  },
  confirm(title: string, content: string, options?: Partial<Omit<ConfirmOptions, 'title' | 'message'>>) {
    const opts: ConfirmOptions = {
      title,
      message: content,
      ...options
    }
    return showConfirm(opts)
  },
}

// 别名 (兼容 Element Plus API)
export const ElMessage = message
export const ElMessageBox = {
  confirm(content: string, title: string, options?: Partial<Omit<ConfirmOptions, 'title' | 'message'>>) {
    return message.confirm(title, content, options).then((confirmed: boolean) => {
      if (!confirmed) {
        return Promise.reject(new Error('cancel'))
      }
      return Promise.resolve()
    })
  },
}

// 全局函数 (兼容旧代码)
declare global {
  interface Window {
    $message: typeof message
  }
}

window.$message = message

export default message
