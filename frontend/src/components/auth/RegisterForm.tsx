import React from 'react'
import { Form, Input, Button, Card, Typography, Space } from 'antd'
import { MailOutlined, LockOutlined, UserOutlined } from '@ant-design/icons'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { registerSchema, type RegisterFormData } from '../../utils/auth-validation'

const { Title, Text } = Typography

interface RegisterFormProps {
  onSubmit: (data: RegisterFormData) => void
  loading?: boolean
}

export const RegisterForm: React.FC<RegisterFormProps> = ({
  onSubmit,
  loading = false,
}) => {
  const {
    control,
    handleSubmit,
    formState: { errors, isValid },
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      email: '',
      password: '',
      confirmPassword: '',
      name: '',
    },
    mode: 'onBlur',
  })

  const handleFormSubmit = (data: RegisterFormData) => {
    onSubmit(data)
  }

  return (
    <Card style={{ maxWidth: 400, width: '100%', padding: '24px' }}>
      <Space direction="vertical" size="large" style={{ width: '100%' }}>
        <div style={{ textAlign: 'center' }}>
          <Title level={2}>Sign Up</Title>
          <Text type="secondary">Create your account to start making predictions!</Text>
        </div>

        <form onSubmit={handleSubmit(handleFormSubmit)}>
          <Space direction="vertical" size="middle" style={{ width: '100%' }}>
            <Controller
              name="name"
              control={control}
              render={({ field }) => (
                <Form.Item
                  validateStatus={errors.name ? 'error' : ''}
                  help={errors.name?.message}
                  style={{ marginBottom: 16 }}
                >
                  <Input
                    {...field}
                    prefix={<UserOutlined />}
                    placeholder="Full Name"
                    autoComplete="name"
                    autoFocus
                    disabled={loading}
                    size="large"
                  />
                </Form.Item>
              )}
            />

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
                  style={{ marginBottom: 16 }}
                >
                  <Input.Password
                    {...field}
                    prefix={<LockOutlined />}
                    placeholder="Password"
                    autoComplete="new-password"
                    disabled={loading}
                    size="large"
                  />
                </Form.Item>
              )}
            />

            <Controller
              name="confirmPassword"
              control={control}
              render={({ field }) => (
                <Form.Item
                  validateStatus={errors.confirmPassword ? 'error' : ''}
                  help={errors.confirmPassword?.message}
                  style={{ marginBottom: 24 }}
                >
                  <Input.Password
                    {...field}
                    prefix={<LockOutlined />}
                    placeholder="Confirm Password"
                    autoComplete="new-password"
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
              {loading ? 'Creating Account...' : 'Create Account'}
            </Button>
          </Space>
        </form>
      </Space>
    </Card>
  )
}

export default RegisterForm
