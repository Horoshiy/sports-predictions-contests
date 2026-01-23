import React, { createContext, useContext } from 'react'
import { message } from 'antd'

interface ToastContextType {
  showToast: (msg: string, severity?: 'success' | 'error' | 'info' | 'warning') => void
}

const ToastContext = createContext<ToastContextType | undefined>(undefined)

export const useToast = () => {
  const context = useContext(ToastContext)
  if (!context) {
    throw new Error('useToast must be used within a ToastProvider')
  }
  return context
}

interface ToastProviderProps {
  children: React.ReactNode
}

export const ToastProvider: React.FC<ToastProviderProps> = ({ children }) => {
  const [messageApi, contextHolder] = message.useMessage()

  const showToast = (msg: string, severity: 'success' | 'error' | 'info' | 'warning' = 'success') => {
    messageApi[severity](msg)
  }

  return (
    <ToastContext.Provider value={{ showToast }}>
      {contextHolder}
      {children}
    </ToastContext.Provider>
  )
}
