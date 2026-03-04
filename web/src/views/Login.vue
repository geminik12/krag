<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import api from '../api'

const router = useRouter()
const message = useMessage()

const form = reactive({ username: '', password: '' })
const isRegister = ref(false)
const loading = ref(false)

const handleLogin = async () => {
  if (!form.username || !form.password) {
    message.warning('请输入用户名和密码')
    return
  }
  loading.value = true
  try {
    if (isRegister.value) {
      await api.post('/v1/users', {
        username: form.username,
        password: form.password,
        nickname: form.username,
        email: form.username + '@example.com',
        phone: '12345678901'
      })
      message.success('注册成功，请登录')
      isRegister.value = false
    } else {
      const res = await api.post('/login', {
        username: form.username,
        password: form.password
      })
      const { token } = res.data
      localStorage.setItem('token', token)
      localStorage.setItem('user', JSON.stringify({ username: form.username }))
      message.success('登录成功')
      router.push('/')
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '操作失败')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-container">
    <n-card style="width: 400px" :title="isRegister ? '注册 Krag' : '登录 Krag'" :bordered="false" size="huge">
      <n-space vertical size="large">
        <n-input v-model:value="form.username" placeholder="用户名" />
        <n-input v-model:value="form.password" type="password" placeholder="密码" @keyup.enter="handleLogin" />
        <n-button type="primary" block :loading="loading" @click="handleLogin">
          {{ isRegister ? '注册' : '登录' }}
        </n-button>
        <div style="text-align: center">
          <n-button text @click="isRegister = !isRegister">
            {{ isRegister ? '已有账号？去登录' : '没有账号？去注册' }}
          </n-button>
        </div>
      </n-space>
    </n-card>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #101014;
}
</style>
