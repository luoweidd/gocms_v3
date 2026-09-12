<template>
  <div class="register-page">
    <div class="register-wrapper">
      <!-- Left Side - Branding -->
      <div class="register-branding">
        <div class="branding-content">
          <div class="brand-logo">
            <span class="logo-icon-big">⚡</span>
            <h1 class="brand-name">GoCMS</h1>
          </div>
          <p class="brand-tagline">现代内容管理系统 - 快速注册，开启管理之旅</p>
          <div class="benefits-list">
            <div class="benefit-item">
              <span class="benefit-icon">✅</span>
              <span>免费开源</span>
            </div>
            <div class="benefit-item">
              <span class="benefit-icon">✅</span>
              <span>多租户支持</span>
            </div>
            <div class="benefit-item">
              <span class="benefit-icon">✅</span>
              <span>完善的权限体系</span>
            </div>
            <div class="benefit-item">
              <span class="benefit-icon">✅</span>
              <span>RESTful API 接口</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Side - Registration Form -->
      <div class="register-form-side">
        <div class="form-container">
          <div class="form-header">
            <h2 class="form-title">创建管理账号</h2>
            <p class="form-subtitle">填写以下信息完成注册</p>
          </div>

          <a-form ref="registerFormRef" :model="registerForm" :rules="registerRules" layout="vertical" class="register-form">
            <!-- Username -->
            <a-form-item label="用户名" name="username">
              <a-input
                v-model:value="registerForm.username"
                placeholder="用户名（3-20位）"
                size="large"
              />
            </a-form-item>

            <!-- Email -->
            <a-form-item label="邮箱" name="email">
              <a-input
                v-model:value="registerForm.email"
                placeholder="电子邮箱"
                size="large"
              />
            </a-form-item>

            <!-- Phone (Optional) -->
            <a-form-item label="手机号" name="phone">
              <a-input
                v-model:value="registerForm.phone"
                placeholder="手机号（可选）"
                size="large"
              />
            </a-form-item>

            <!-- Password -->
            <a-form-item label="密码" name="password">
              <a-input-password
                v-model:value="registerForm.password"
                placeholder="密码（至少6位）"
                size="large"
              />
            </a-form-item>

            <!-- Confirm Password -->
            <a-form-item label="确认密码" name="confirmPassword">
              <a-input-password
                v-model:value="registerForm.confirmPassword"
                placeholder="确认密码"
                size="large"
              />
            </a-form-item>

            <!-- Terms -->
            <div class="terms-check">
              <a-checkbox v-model:checked="registerForm.acceptTerms">
                我已阅读并同意<a href="#" class="terms-link">服务条款</a>和<a href="#" class="terms-link">隐私政策</a>
              </a-checkbox>
            </div>

            <!-- Submit Button -->
            <a-button
              type="primary"
              size="large"
              class="register-btn"
              :loading="loading"
              @click="handleRegister"
            >
              {{ loading ? '注册中...' : '注 册' }}
            </a-button>

            <!-- Login Link -->
            <div class="form-footer">
              <span>已有账号？</span>
              <router-link to="/login" class="login-link">立即登录</router-link>
            </div>
          </a-form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { message as antMessage } from 'ant-design-vue'

const router = useRouter()
const registerFormRef = ref(null)
const loading = ref(false)

const validateConfirmPassword = (_rule: unknown, value: string) => {
  if (value !== registerForm.password) {
    return Promise.reject('两次输入的密码不一致')
  }
  return Promise.resolve()
}

const registerForm = reactive({
  username: '',
  email: '',
  phone: '',
  password: '',
  confirmPassword: '',
  acceptTerms: false
})

const registerRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度为3-20位', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少为6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const handleRegister = async () => {
  if (!registerFormRef.value) return
  
  try {
    await (registerFormRef.value as any).validate()
    if (!registerForm.acceptTerms) {
      antMessage.warning('请先阅读并同意服务条款')
      return
    }
    loading.value = true
    try {
      // TODO: 调用注册 API
      // await register(registerForm)
      antMessage.success('注册成功，请前往登录')
      router.push('/login')
    } finally {
      loading.value = false
    }
  } catch {
    antMessage.error('注册失败，请重试')
  }
}
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.register-wrapper {
  display: flex;
  width: 900px;
  max-width: 95vw;
  background: #fff;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

/* Left Branding Side */
.register-branding {
  flex: 1;
  background: linear-gradient(135deg, #4F46E5, #7C3AED);
  padding: 60px 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.branding-content {
  text-align: center;
}

.logo-icon-big {
  font-size: 48px;
}

.brand-name {
  font-size: 32px;
  font-weight: 700;
  margin-left: 8px;
}

.brand-tagline {
  font-size: 15px;
  opacity: 0.9;
  margin: 16px 0 40px;
  line-height: 1.6;
}

.benefits-list {
  text-align: left;
  max-width: 280px;
  margin: 0 auto;
}

.benefit-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 0;
  font-size: 15px;
  opacity: 0.9;
}

.benefit-icon {
  font-size: 18px;
}

/* Right Form Side */
.register-form-side {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
  background: #fff;
}

.form-container {
  width: 100%;
  max-width: 360px;
}

.form-header {
  margin-bottom: 32px;
  text-align: center;
}

.form-title {
  font-size: 24px;
  font-weight: 700;
  color: #1F2937;
  margin-bottom: 8px;
}

.form-subtitle {
  font-size: 14px;
  color: #6B7280;
}

.register-form {
  margin-top: 24px;
}

.terms-check {
  margin-bottom: 24px;
}

.terms-link {
  color: #4F46E5;
  text-decoration: none;
}

.terms-link:hover {
  text-decoration: underline;
}

.register-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 8px;
  background: linear-gradient(135deg, #4F46E5, #7C3AED);
  border: none;
}

.register-btn:hover {
  opacity: 0.9;
}

.form-footer {
  text-align: center;
  margin-top: 24px;
  font-size: 14px;
  color: #6B7280;
}

.login-link {
  color: #4F46E5;
  text-decoration: none;
  font-weight: 600;
}

.login-link:hover {
  text-decoration: underline;
}

/* Responsive */
@media (max-width: 768px) {
  .register-wrapper {
    flex-direction: column;
    width: 95vw;
  }

  .register-branding {
    padding: 40px 20px;
  }

  .brand-name {
    font-size: 24px;
  }

  .brand-tagline {
    margin-bottom: 24px;
  }

  .register-form-side {
    padding: 30px 20px;
  }
}
</style>