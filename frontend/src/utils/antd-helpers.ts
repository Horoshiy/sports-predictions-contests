import { message, notification } from 'antd';
import type { ArgsProps as MessageArgsProps } from 'antd/es/message';
import type { ArgsProps as NotificationArgsProps } from 'antd/es/notification';

/**
 * Show success message
 */
export const showSuccess = (content: string, duration = 3) => {
  message.success(content, duration);
};

/**
 * Show error message
 */
export const showError = (content: string, duration = 3) => {
  message.error(content, duration);
};

/**
 * Show warning message
 */
export const showWarning = (content: string, duration = 3) => {
  message.warning(content, duration);
};

/**
 * Show info message
 */
export const showInfo = (content: string, duration = 3) => {
  message.info(content, duration);
};

/**
 * Show loading message
 */
export const showLoading = (content: string) => {
  return message.loading(content, 0);
};

/**
 * Show success notification with title and description
 */
export const notifySuccess = (title: string, description?: string) => {
  notification.success({
    message: title,
    description,
    placement: 'topRight',
  });
};

/**
 * Show error notification with title and description
 */
export const notifyError = (title: string, description?: string) => {
  notification.error({
    message: title,
    description,
    placement: 'topRight',
  });
};

/**
 * Show warning notification with title and description
 */
export const notifyWarning = (title: string, description?: string) => {
  notification.warning({
    message: title,
    description,
    placement: 'topRight',
  });
};

/**
 * Show info notification with title and description
 */
export const notifyInfo = (title: string, description?: string) => {
  notification.info({
    message: title,
    description,
    placement: 'topRight',
  });
};

/**
 * Destroy all messages
 */
export const destroyAllMessages = () => {
  message.destroy();
};

/**
 * Destroy all notifications
 */
export const destroyAllNotifications = () => {
  notification.destroy();
};
