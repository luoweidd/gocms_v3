/// <reference types="vite/client" />
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/dashboard'
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/pages/login/LoginPage.vue'),
      meta: { title: '登录' }
    },
    {
      path: '/',
      component: () => import('@/layouts/MainLayout.vue'),
      redirect: '/dashboard',
      children: [
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('@/pages/dashboard/DashboardPage.vue'),
          meta: { title: '仪表盘' }
        },
        {
          path: 'articles',
          redirect: '/articles/list',
          children: [
            {
              path: 'list',
              name: 'article-list',
              component: () => import('@/pages/articles/ArticleListPage.vue'),
              meta: { title: '文章列表' }
            },
            {
              path: 'create',
              name: 'article-create',
              component: () => import('@/pages/articles/ArticleFormPage.vue'),
              meta: { title: '发布文章' }
            },
            {
              path: ':id',
              name: 'article-detail',
              component: () => import('@/pages/articles/ArticleDetailPage.vue'),
              meta: { title: '文章详情' }
            },
            {
              path: 'edit/:id',
              name: 'article-edit',
              component: () => import('@/pages/articles/ArticleFormPage.vue'),
              meta: { title: '编辑文章' }
            },
            {
              path: 'categories',
              name: 'article-categories',
              component: () => import('@/pages/articles/CategoryPage.vue'),
              meta: { title: '分类管理' }
            }
          ]
        },
        {
          path: 'videos',
          redirect: '/videos/list',
          children: [
            {
              path: 'list',
              name: 'video-list',
              component: () => import('@/pages/videos/VideoListPage.vue'),
              meta: { title: '视频列表' }
            },
            {
              path: 'create',
              name: 'video-create',
              component: () => import('@/pages/videos/VideoFormPage.vue'),
              meta: { title: '上传视频' }
            },
            {
              path: 'edit/:id',
              name: 'video-edit',
              component: () => import('@/pages/videos/VideoFormPage.vue'),
              meta: { title: '编辑视频' }
            },
            {
              path: 'categories',
              name: 'video-categories',
              component: () => import('@/pages/videos/CategoryPage.vue'),
              meta: { title: '分类管理' }
            }
          ]
        },
        {
          path: 'comments',
          redirect: '/comments/list',
          children: [
            {
              path: 'list',
              name: 'comment-list',
              component: () => import('@/pages/comments/CommentListPage.vue'),
              meta: { title: '评论列表' }
            }
          ]
        },
        {
          path: 'users',
          redirect: '/users/list',
          children: [
            {
              path: 'list',
              name: 'user-list',
              component: () => import('@/pages/users/UserListPage.vue'),
              meta: { title: '用户列表' }
            },
            {
              path: 'profile',
              name: 'user-profile',
              component: () => import('@/pages/users/ProfilePage.vue'),
              meta: { title: '个人信息' }
            },
            {
              path: 'role',
              name: 'user-role',
              component: () => import('@/pages/system/RolePage.vue'),
              meta: { title: '角色管理' }
            },
            {
              path: 'permission',
              name: 'user-permission',
              component: () => import('@/pages/system/PermissionPage.vue'),
              meta: { title: '权限管理' }
            }
          ]
        },
        {
          path: 'message',
          redirect: '/message/notifications',
          children: [
            {
              path: 'notifications',
              name: 'message-notifications',
              component: () => import('@/pages/messages/NotificationPage.vue'),
              meta: { title: '我的通知' }
            },
            {
              path: 'system',
              name: 'message-system',
              component: () => import('@/pages/messages/SystemMessagePage.vue'),
              meta: { title: '系统消息管理' }
            },
            {
              path: 'categories',
              name: 'message-categories',
              component: () => import('@/pages/messages/CategoryPage.vue'),
              meta: { title: '消息类别管理' }
            }
          ]
        },
        {
          path: 'audit',
          redirect: '/audit/workbench',
          children: [
            {
              path: 'workbench',
              name: 'audit-workbench',
              component: () => import('@/pages/audit/index.vue'),
              meta: { title: '审核工作台' }
            },
            {
              path: 'articles',
              name: 'audit-articles',
              component: () => import('@/pages/audit/ArticleAuditPage.vue'),
              meta: { title: '文章审核' }
            },
            {
              path: 'videos',
              name: 'audit-videos',
              component: () => import('@/pages/audit/VideoAuditPage.vue'),
              meta: { title: '视频审核' }
            }
          ]
        },
        {
          path: 'media',
          redirect: '/media/dashboard',
          children: [
            {
              path: 'dashboard',
              name: 'media-dashboard',
              component: () => import('@/pages/media/Dashboard.vue'),
              meta: { title: '媒体中心' }
            },
            {
              path: 'images',
              name: 'media-images',
              component: () => import('@/pages/media/ImageManager.vue'),
              meta: { title: '图片管理' }
            },
            {
              path: 'files',
              name: 'media-files',
              component: () => import('@/pages/media/FileManager.vue'),
              meta: { title: '文件管理' }
            },
            {
              path: 'audio',
              name: 'media-audio',
              component: () => import('@/pages/media/AudioManager.vue'),
              meta: { title: '音频管理' }
            }
          ]
        },
        {
          path: 'system',
          redirect: '/system/menus',
          children: [
            {
              path: 'menus',
              name: 'system-menus',
              component: () => import('@/pages/system/MenuPage.vue'),
              meta: { title: '菜单管理' }
            },
            {
              path: 'roles',
              name: 'system-roles',
              component: () => import('@/pages/system/RolePage.vue'),
              meta: { title: '角色权限' }
            },
            {
              path: 'logs',
              name: 'system-logs',
              component: () => import('@/pages/system/OperationLogPage.vue'),
              meta: { title: '操作日志' }
            },
            {
              path: 'backup',
              name: 'system-backup',
              component: () => import('@/pages/system/backup/index.vue'),
              meta: { title: '数据备份' }
            },
            {
              path: 'seo',
              name: 'system-seo',
              component: () => import('@/pages/system/SeoPage.vue'),
              meta: { title: 'SEO管理' }
            },
            {
              path: 'config',
              name: 'system-config',
              component: () => import('@/pages/system/SystemConfigPage.vue'),
              meta: { title: '系统设置' }
            }
          ]
        },
        {
          path: 'ngac',
          redirect: '/ngac/permissions',
          children: [
            {
              path: 'permissions',
              name: 'ngac-permissions',
              component: () => import('@/pages/system/PermissionPage.vue'),
              meta: { title: '权限列表' }
            }
          ]
        },
        {
          path: 'tenants',
          redirect: '/tenants/list',
          children: [
            {
              path: 'list',
              name: 'tenant-list',
              component: () => import('@/pages/system/tenant/index.vue'),
              meta: { title: '租户列表' }
            }
          ]
        },
        {
          path: 'cron',
          redirect: '/cron/tasks',
          children: [
            {
              path: 'tasks',
              name: 'cron-tasks',
              component: () => import('@/pages/cron/TaskPage.vue'),
              meta: { title: '任务列表' }
            },
            {
              path: 'logs',
              name: 'cron-logs',
              component: () => import('@/pages/cron/LogPage.vue'),
              meta: { title: '执行日志' }
            }
          ]
        }
      ]
    }
  ]
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = (to.meta.title as string) || 'GoCMS 管理系统'

  // 获取 token
  const token = localStorage.getItem('token')

  // 无需登录的页面
  const publicPages = ['/login']

  if (!token && !publicPages.includes(to.path)) {
    // 未登录且访问需要认证的页面，重定向到登录页
    next(`/login?redirect=${to.path}`)
  } else if (token && to.path === '/login') {
    // 已登录且访问登录页，重定向到首页
    next('/dashboard')
  } else {
    next()
  }
})

export default router