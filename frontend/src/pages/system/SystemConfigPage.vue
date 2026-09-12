<template>
  <div class="page-container">
    <!-- Page Header -->
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 mb-1">系统设置</h1>
        <p class="text-sm text-gray-500">管理系统基本配置和参数</p>
      </div>
    </div>

    <!-- Settings Tabs -->
    <div class="ta-card p-6">
      <div class="border-b border-gray-200 mb-6">
        <nav class="-mb-px flex gap-8">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            class="pb-4 px-1 border-b-2 font-medium text-sm transition-colors"
            :class="
              activeTab === tab.id
                ? 'border-primary text-primary'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            "
          >
            {{ tab.name }}
          </button>
        </nav>
      </div>

      <!-- 基本设置 -->
      <div v-if="activeTab === 'basic'" class="space-y-6">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div class="form-group">
            <label class="block text-sm font-medium text-gray-700 mb-2">网站名称</label>
            <input
              type="text"
              v-model="config.site_name"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-colors"
              placeholder="请输入网站名称"
            />
          </div>
          <div class="form-group">
            <label class="block text-sm font-medium text-gray-700 mb-2">网站URL</label>
            <input
              type="url"
              v-model="config.site_url"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-colors"
              placeholder="https://example.com"
            />
          </div>
          <div class="form-group">
            <label class="block text-sm font-medium text-gray-700 mb-2">网站描述</label>
            <textarea
              v-model="config.site_description"
              rows="3"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-colors resize-none"
              placeholder="请输入网站描述"
            ></textarea>
          </div>
          <div class="form-group">
            <label class="block text-sm font-medium text-gray-700 mb-2">站点图标 (Favicon)</label>
            <input
              type="text"
              v-model="config.favicon_url"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-colors"
              placeholder="图标URL地址"
            />
          </div>
        </div>
        <div class="flex justify-end pt-4">
          <button
            @click="saveConfig('basic')"
            class="px-6 py-2 bg-primary text-white rounded-lg hover:bg-primary/90 transition-colors flex items-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            保存设置
          </button>
        </div>
      </div>

      <!-- 安全设置 -->
      <div v-if="activeTab === 'security'" class="space-y-6">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div class="form-group">
            <label class="block text-sm font-medium text-gray-700 mb-2">密码最小长度</label>
            <input
              type="number"
              v-model="config.min_password_length"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-colors"
            />
          </div>
          <div class="form-group">
            <label class="block text-sm font-medium text-gray-700 mb-2">登录最大尝试次数</label>
            <input
              type="number"
              v-model="config.max_login_attempts"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-colors"
            />
          </div>
          <div class="form-group">
            <label class="block text-sm font-medium text-gray-700 mb-2">会话过期时间（分钟）</label>
            <input
              type="number"
              v-model="config.session_timeout"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-colors"
            />
          </div>
          <div class="form-group">
            <label class="block text-sm font-medium text-gray-700 mb-2">验证码启用状态</label>
            <select
              v-model="config.captcha_enabled"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-colors"
            >
              <option :value="true">启用</option>
              <option :value="false">禁用</option>
            </select>
          </div>
        </div>
        <!-- 安全开关 -->
        <div class="pt-4 border-t border-gray-200">
          <h3 class="text-sm font-semibold text-gray-900 mb-4">功能开关</h3>
          <div class="space-y-4">
            <div class="flex items-center justify-between py-3">
              <div>
                <p class="text-sm font-medium text-gray-700">双因素认证</p>
                <p class="text-xs text-gray-500">启用后用户登录需要额外验证</p>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="config.two_factor_enabled" class="sr-only peer" />
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary/20 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
              </label>
            </div>
            <div class="flex items-center justify-between py-3">
              <div>
                <p class="text-sm font-medium text-gray-700">API访问限制</p>
                <p class="text-xs text-gray-500">限制每个IP的API请求频率</p>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="config.api_rate_limit" class="sr-only peer" />
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary/20 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
              </label>
            </div>
          </div>
        </div>
        <div class="flex justify-end pt-4">
          <button
            @click="saveConfig('security')"
            class="px-6 py-2 bg-primary text-white rounded-lg hover:bg-primary/90 transition-colors flex items-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            保存设置
          </button>
        </div>
      </div>

      <!-- 存储设置 -->
      <div v-if="activeTab === 'storage'" class="space-y-6">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div class="form-group">
            <label class="block text-sm font-medium text-gray-700 mb-2">存储类型</label>
            <select
              v-model="config.storage_type"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-colors"
            >
              <option value="local">本地存储</option>
              <option value="oss">对象存储 (OSS)</option>
              <option value="cdn">CDN加速</option>
            </select>
          </div>
          <div class="form-group">
            <label class="block text-sm font-medium text-gray-700 mb-2">最大上传文件大小（MB）</label>
            <input
              type="number"
              v-model="config.max_file_size"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-colors"
            />
          </div>
          <div class="form-group lg:col-span-2">
            <label class="block text-sm font-medium text-gray-700 mb-2">允许的上传文件类型</label>
            <input
              type="text"
              v-model="config.allowed_file_types"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-colors"
              placeholder="例如: jpg,png,gif,mp4,pdf"
            />
            <p class="text-xs text-gray-500 mt-1">多个类型用逗号分隔</p>
          </div>
        </div>
        <div class="flex justify-end pt-4">
          <button
            @click="saveConfig('storage')"
            class="px-6 py-2 bg-primary text-white rounded-lg hover:bg-primary/90 transition-colors flex items-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            保存设置
          </button>
        </div>
      </div>

      <!-- 邮件设置 -->
      <div v-if="activeTab === 'email'" class="space-y-6">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div class="form-group">
            <label class="block text-sm font-medium text-gray-700 mb-2">SMTP服务器</label>
            <input
              type="text"
              v-model="config.smtp_host"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-colors"
              placeholder="smtp.example.com"
            />
          </div>
          <div class="form-group">
            <label class="block text-sm font-medium text-gray-700 mb-2">SMTP端口</label>
            <input
              type="number"
              v-model="config.smtp_port"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-colors"
            />
          </div>
          <div class="form-group">
            <label class="block text-sm font-medium text-gray-700 mb-2">发件人邮箱</label>
            <input
              type="email"
              v-model="config.smtp_user"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-colors"
            />
          </div>
          <div class="form-group">
            <label class="block text-sm font-medium text-gray-700 mb-2">发件人密码</label>
            <input
              type="password"
              v-model="config.smtp_pass"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-colors"
            />
          </div>
          <div class="form-group lg:col-span-2">
            <label class="block text-sm font-medium text-gray-700 mb-2">发件人名称</label>
            <input
              type="text"
              v-model="config.from_name"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-colors"
            />
          </div>
        </div>
        <div class="flex justify-end pt-4 gap-3">
          <button
            @click="testEmail"
            class="px-6 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors flex items-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
            测试发送
          </button>
          <button
            @click="saveConfig('email')"
            class="px-6 py-2 bg-primary text-white rounded-lg hover:bg-primary/90 transition-colors flex items-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            保存设置
          </button>
        </div>
      </div>

      <!-- SEO设置 -->
      <div v-if="activeTab === 'seo'" class="space-y-6">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div class="form-group lg:col-span-2">
            <label class="block text-sm font-medium text-gray-700 mb-2">SEO标题</label>
            <input
              type="text"
              v-model="config.seo_title"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-colors"
              placeholder="网站标题"
            />
          </div>
          <div class="form-group lg:col-span-2">
            <label class="block text-sm font-medium text-gray-700 mb-2">SEO描述</label>
            <textarea
              v-model="config.seo_description"
              rows="3"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-colors resize-none"
              placeholder="网站描述"
            ></textarea>
          </div>
          <div class="form-group lg:col-span-2">
            <label class="block text-sm font-medium text-gray-700 mb-2">SEO关键词</label>
            <input
              type="text"
              v-model="config.seo_keywords"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-colors"
              placeholder="多个关键词用逗号分隔"
            />
          </div>
        </div>
        <div class="flex justify-end pt-4">
          <button
            @click="saveConfig('seo')"
            class="px-6 py-2 bg-primary text-white rounded-lg hover:bg-primary/90 transition-colors flex items-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            保存设置
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { message } from '@/utils/message'

const activeTab = ref('basic')

const tabs = [
  { id: 'basic', name: '基本设置' },
  { id: 'security', name: '安全设置' },
  { id: 'storage', name: '存储设置' },
  { id: 'email', name: '邮件设置' },
  { id: 'seo', name: 'SEO设置' },
]

const config = reactive({
  site_name: '',
  site_url: '',
  site_description: '',
  favicon_url: '',
  min_password_length: 8,
  max_login_attempts: 5,
  session_timeout: 720,
  captcha_enabled: true,
  two_factor_enabled: false,
  api_rate_limit: false,
  storage_type: 'local',
  max_file_size: 100,
  allowed_file_types: 'jpg,png,gif,mp4,pdf,doc,docx,xlsx,xls',
  smtp_host: '',
  smtp_port: 587,
  smtp_user: '',
  smtp_pass: '',
  from_name: '',
  seo_title: '',
  seo_description: '',
  seo_keywords: '',
})

const loadConfig = async () => {
  // TODO: 从后端加载系统配置
  console.log('加载系统配置...')
}

const saveConfig = async (tab: string) => {
  try {
    // TODO: 实现保存逻辑，调用后端 API
    console.log('保存配置:', tab, config)
    message.success('设置已保存')
  } catch (error) {
    message.error('保存失败，请重试')
  }
}

const testEmail = async () => {
  // TODO: 实现测试发送功能
  console.log('测试邮件发送:', config.smtp_host, config.smtp_port, config.smtp_user)
}

loadConfig()
</script>

<style scoped>
.page-container {
  max-width: 1200px;
  margin: 0 auto;
}
</style>