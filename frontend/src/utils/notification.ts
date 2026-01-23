import { message } from 'antd'

type NotificationType = 'success' | 'error' | 'info' | 'warning'

export const showNotification = (content: string, type: NotificationType = 'info') => {
  message[type](content)
}

export const showSuccess = (content: string) => showNotification(content, 'success')
export const showError = (content: string) => showNotification(content, 'error')
export const showInfo = (content: string) => showNotification(content, 'info')
export const showWarning = (content: string) => showNotification(content, 'warning')
