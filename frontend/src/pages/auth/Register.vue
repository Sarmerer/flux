<script setup lang="ts">
import { Eye, EyeOff, Loader2 } from 'lucide-vue-next'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

import { ROUTE_PATHS } from '@/constants/routes'
import { useAuthStore } from '@/stores/auth'

import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

const router = useRouter()
const authStore = useAuthStore()

const name = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const showPassword = ref(false)
const showConfirmPassword = ref(false)
const isLoading = ref(false)

const passwordsMatch = ref(true)

const handleSubmit = async (e: Event) => {
  e.preventDefault()

  if (!name.value || !email.value || !password.value || !confirmPassword.value) return

  if (password.value !== confirmPassword.value) {
    passwordsMatch.value = false
    return
  }

  passwordsMatch.value = true
  isLoading.value = true

  try {
    await authStore.register(name.value, email.value, password.value)
    router.push(ROUTE_PATHS.HOME)
  } catch (error) {
  } finally {
    isLoading.value = false
  }
}

const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value
}

const toggleConfirmPasswordVisibility = () => {
  showConfirmPassword.value = !showConfirmPassword.value
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div class="text-center">
        <h1 class="text-3xl font-bold text-gray-900">FlowDB</h1>
        <p class="mt-2 text-sm text-gray-600">Create your account</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Sign Up</CardTitle>
          <CardDescription> Create a new account to get started with FlowDB </CardDescription>
        </CardHeader>
        <CardContent>
          <form @submit="handleSubmit" class="space-y-4">
            <div class="space-y-2">
              <Label for="name">Full Name</Label>
              <Input
                id="name"
                v-model="name"
                type="text"
                placeholder="Enter your full name"
                required
                :disabled="isLoading"
              />
            </div>

            <div class="space-y-2">
              <Label for="email">Email</Label>
              <Input
                id="email"
                v-model="email"
                type="email"
                placeholder="Enter your email"
                required
                :disabled="isLoading"
              />
            </div>

            <div class="space-y-2">
              <Label for="password">Password</Label>
              <div class="relative">
                <Input
                  id="password"
                  v-model="password"
                  :type="showPassword ? 'text' : 'password'"
                  placeholder="Enter your password"
                  required
                  :disabled="isLoading"
                />
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  class="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                  @click="togglePasswordVisibility"
                  :disabled="isLoading"
                >
                  <Eye v-if="!showPassword" class="h-4 w-4" />
                  <EyeOff v-else class="h-4 w-4" />
                </Button>
              </div>
            </div>

            <div class="space-y-2">
              <Label for="confirmPassword">Confirm Password</Label>
              <div class="relative">
                <Input
                  id="confirmPassword"
                  v-model="confirmPassword"
                  :type="showConfirmPassword ? 'text' : 'password'"
                  placeholder="Confirm your password"
                  required
                  :disabled="isLoading"
                />
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  class="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                  @click="toggleConfirmPasswordVisibility"
                  :disabled="isLoading"
                >
                  <Eye v-if="!showConfirmPassword" class="h-4 w-4" />
                  <EyeOff v-else class="h-4 w-4" />
                </Button>
              </div>
            </div>

            <Alert v-if="!passwordsMatch" variant="destructive">
              <AlertDescription>Passwords do not match</AlertDescription>
            </Alert>

            <Alert v-if="authStore.error" variant="destructive">
              <AlertDescription>{{ authStore.error }}</AlertDescription>
            </Alert>

            <Button type="submit" class="w-full" :disabled="isLoading">
              <Loader2 v-if="isLoading" class="mr-2 h-4 w-4 animate-spin" />
              {{ isLoading ? 'Creating account...' : 'Create Account' }}
            </Button>
          </form>

          <div class="mt-6 text-center">
            <p class="text-sm text-gray-600">
              Already have an account?
              <router-link :to="ROUTE_PATHS.LOGIN" class="font-medium text-blue-600 hover:text-blue-500">
                Sign in
              </router-link>
            </p>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
