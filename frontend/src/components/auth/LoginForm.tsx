import React from 'react'
import { Form, Input, Button, Card, Typography, Space } from 'antd'
import { MailOutlined, LockOutlined } from '@ant-design/icons'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { loginSchema, type LoginFormData } from '../../utils/auth-validation'

const { Title, Text } = Typography

interface LoginFormProps {
  onSubmit: (data: LoginFormData) => void
  loading?: boolean
}

export const LoginForm: React.FC<LoginFormProps> = ({
  onSubmit,
  loading = false,
}) => {
  const {
    control,
    handleSubmit,
    formState: { errors, isValid },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: '',
      password: '',
    },
    mode: 'onBlur',
  })

  const handleFormSubmit = (data: LoginFormData) => {
    onSubmit(data)
  }

  return (
    <Card style={{ maxWidth: 400, width: '100%', padding: '24px' }}>
      <Space direction="vertical" size="large" style={{ width: '100%' }}>
        <div style={{ textAlign: 'center' }}>
          <Title level={2}>Sign In</Title>
          <Text type="secondary">Welcome back! Please sign in to your account.</Text>
        </div>

        <form onSubmit={handleSubmit(handleFormSubmit)}>
          <Space direction="vertical" size="middle" style={{ width: '100%' }}>
            <Controller
              name="email"
              control={control}
              render={({ field }) => (
                <Form.Item
                  validateStatus={errors.email ? 'error' : ''}
                  help={errors.email?.message}
                  style={{ marginBottom: 16 }}
                >
                  <Input
                    {...field}
                    prefix={<MailOutlined />}
                    placeholder="Email Address"
                    type="email"
                    autoComplete="email"
                    autoFocus
                    disabled={loading}
                    size="large"
                  />
                </Form.Item>
              )}
            />

            <Controller
              name="password"
              control={control}
              render={({ field }) => (
                <Form.Item
                  validateStatus={errors.password ? 'error' : ''}
                  help={errors.password?.message}
                  style={{ marginBottom: 24 }}
                >
                  <Input.Password
                    {...field}
                    prefix={<LockOutlined />}
                    placeholder="Password"
                    autoComplete="current-password"
                    disabled={loading}
                    size="large"
                  />
                </Form.Item>
              )}
            />

            <Button
              type="primary"
              htmlType="submit"
              block
              size="large"
              disabled={!isValid || loading}
              loading={loading}
            >
              {loading ? 'Signing In...' : 'Sign In'}
            </Button>
          </Space>
        </form>
      </Space>
    </Card>
  )
}

export default LoginForm
