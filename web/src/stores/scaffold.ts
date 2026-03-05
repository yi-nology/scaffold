import { defineStore } from 'pinia'
import { ref } from 'vue'

interface Variable {
  name: string
  type: 'string' | 'boolean' | 'enum' | 'number'
  default: any
  prompt: string
  options?: string[]
  required?: boolean
  description?: string
  group?: string
}

interface TemplateMeta {
  id: string
  name: string
  description: string
  author: string
  tags: string[]
  version: string
  repository?: string
}

interface TagInfo {
  name: string
  message: string
}

export const useScaffoldStore = defineStore('scaffold', () => {
  const templates = ref<TemplateMeta[]>([])
  const selectedTemplate = ref<TemplateMeta | null>(null)
  const variables = ref<Variable[]>([])
  const values = ref<Record<string, any>>({})
  const zipBlob = ref<Blob | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const availableTags = ref<TagInfo[]>([])
  const selectedVersion = ref<string>('')

  const API_BASE = '/api'

  async function fetchTemplates() {
    loading.value = true
    error.value = null
    try {
      const res = await fetch(`${API_BASE}/templates`)
      const response = await res.json()
      if (response.code === 0) {
        templates.value = response.data || []
      } else {
        error.value = response.message || 'Failed to fetch templates'
      }
    } catch (e: any) {
      error.value = e.message
    } finally {
      loading.value = false
    }
  }

  async function selectTemplate(id: string) {
    loading.value = true
    error.value = null
    try {
      const res = await fetch(`${API_BASE}/templates/${id}`)
      const response = await res.json()
      if (response.code === 0) {
        const data = response.data
        const config = data.config || {}
        selectedTemplate.value = {
          id: data.id || id,
          name: config.name || '',
          description: config.description || '',
          author: config.author || '',
          tags: config.tags || [],
          version: config.version || '',
          repository: data.repository,
        }
        variables.value = config.variables || []
        
        values.value = {}
        for (const v of variables.value) {
          values.value[v.name] = v.default
        }
        
        // 获取远端模板的版本 tags
        await fetchTemplateTags(id)
      } else {
        error.value = response.message || 'Failed to fetch template'
      }
    } catch (e: any) {
      error.value = e.message
    } finally {
      loading.value = false
    }
  }

  async function fetchTemplateTags(id: string) {
    try {
      const res = await fetch(`${API_BASE}/templates/${id}/tags`)
      const response = await res.json()
      if (response.code === 0) {
        availableTags.value = response.data || []
        selectedVersion.value = ''
      } else {
        availableTags.value = []
      }
    } catch (e: any) {
      availableTags.value = []
    }
  }

  async function generateProject() {
    if (!selectedTemplate.value) return
    
    loading.value = true
    error.value = null
    try {
      const res = await fetch(`${API_BASE}/generate`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          template_id: selectedTemplate.value.id,
          variables: values.value,
          version: selectedVersion.value || undefined
        })
      })

      if (!res.ok) {
        const errData = await res.json()
        throw new Error(errData.message || 'Generation failed')
      }

      zipBlob.value = await res.blob()
    } catch (e: any) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  function downloadZip() {
    if (!zipBlob.value) return
    
    const projectName = values.value.project_name || 'project'
    const url = URL.createObjectURL(zipBlob.value)
    const a = document.createElement('a')
    a.href = url
    a.download = `${projectName}.zip`
    a.click()
    URL.revokeObjectURL(url)
  }

  function setValue(name: string, value: any) {
    values.value[name] = value
  }

  function reset() {
    selectedTemplate.value = null
    variables.value = []
    values.value = {}
    zipBlob.value = null
    error.value = null
    availableTags.value = []
    selectedVersion.value = ''
  }

  function setVersion(version: string) {
    selectedVersion.value = version
  }

  async function addTemplate(id: string, repoUrl: string, accessKey?: string): Promise<{taskId: string}> {
    loading.value = true
    error.value = null
    try {
      const headers: Record<string, string> = { 'Content-Type': 'application/json' }
      if (accessKey) {
        headers['Authorization'] = `Bearer ${accessKey}`
      }
      
      const res = await fetch(`${API_BASE}/templates`, {
        method: 'POST',
        headers,
        body: JSON.stringify({ id, repository: repoUrl })
      })
      const response = await res.json()
      if (response.code !== 0) {
        throw new Error(response.message || 'Failed to add template')
      }
      
      // 返回任务ID用于后续查询状态
      return { taskId: response.data.task_id }
    } catch (e: any) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function getTaskStatus(taskId: string): Promise<any> {
    try {
      const res = await fetch(`${API_BASE}/tasks/${taskId}`)
      const response = await res.json()
      if (response.code !== 0) {
        throw new Error(response.message || 'Failed to get task status')
      }
      return response.data
    } catch (e: any) {
      throw e
    }
  }

  async function getAllTasks(): Promise<any[]> {
    try {
      const res = await fetch(`${API_BASE}/tasks`)
      const response = await res.json()
      if (response.code !== 0) {
        throw new Error(response.message || 'Failed to get tasks')
      }
      return response.data
    } catch (e: any) {
      throw e
    }
  }

  async function deleteTemplate(id: string, accessKey?: string): Promise<void> {
    loading.value = true
    error.value = null
    try {
      const headers: Record<string, string> = { 'Content-Type': 'application/json' }
      if (accessKey) {
        headers['Authorization'] = `Bearer ${accessKey}`
      }
      
      const res = await fetch(`${API_BASE}/templates/delete`, {
        method: 'POST',
        headers,
        body: JSON.stringify({ id })
      })
      const response = await res.json()
      if (response.code !== 0) {
        throw new Error(response.message || 'Failed to delete template')
      }
      await fetchTemplates()
    } catch (e: any) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  return {
    templates,
    selectedTemplate,
    variables,
    values,
    zipBlob,
    loading,
    error,
    availableTags,
    selectedVersion,
    fetchTemplates,
    selectTemplate,
    generateProject,
    downloadZip,
    setValue,
    setVersion,
    reset,
    addTemplate,
    deleteTemplate,
    getTaskStatus,
    getAllTasks
  }
})
