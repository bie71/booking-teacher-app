<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white p-6 rounded-lg w-full max-w-lg mx-4 shadow-lg">
      <h2 class="text-xl font-bold mb-4">Edit Teacher</h2>

      <form @submit.prevent="submitForm">
        <!-- Name -->
        <div class="mb-4">
          <label class="block text-sm font-medium mb-1">Name</label>
          <input v-model="form.name" type="text" class="input" />
        </div>

        <!-- Price -->
        <div class="mb-4">
          <label class="block text-sm font-medium mb-1">Price per Hour</label>
          <input v-model.number="form.price_per_hour" type="number" class="input" />
        </div>

        <!-- Bio -->
        <div class="mb-4">
          <label for="bio" class="block text-sm font-medium text-gray-700">Bio</label>
          <textarea 
            v-model="form.bio"
            id="bio"
            rows="3"
            class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          ></textarea>
        </div>

        <!-- Start Time -->
        <div class="mb-4">
          <label class="block text-sm font-medium mb-1">Available Start Time</label>
          <input v-model="form.available_start_time" type="time" class="input" />
        </div>

        <!-- End Time -->
        <div class="mb-4">
          <label class="block text-sm font-medium mb-1">Available End Time</label>
          <input v-model="form.available_end_time" type="time" class="input" />
        </div>

        <!-- Profile Image -->
        <div class="bg-white shadow rounded-lg p-6 mb-4">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Profile Image</h3>
          <div>
            <label for="profile_image" class="block text-sm font-medium text-gray-700">Profile Image</label>
            <div class="mt-1 flex items-center">
              <div class="mr-4">
                <img 
                  :src="previewImage || form.profile_image" 
                  alt="Preview" 
                  class="h-20 w-20 rounded-full object-cover"
                />
              </div>
              <input 
                type="file"
                id="profile_image"
                accept="image/*"
                @change="handleImageUpload"
                class="mt-1 block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-medium file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
              />
            </div>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex justify-end space-x-2">
          <button type="button" @click="$emit('close')" class="btn-secondary">Cancel</button>
          <button type="submit" class="btn-primary">Save</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { teacherService } from '@/services/api'
import { useUIStore } from '@/stores/ui'
import { useTeachersStore } from '@/stores/teachers'

const teachersStore = useTeachersStore()
const uiStore = useUIStore()

const props = defineProps({ teacher: Object })
const emit = defineEmits(['close'])

const previewImage = ref('')
const form = ref({ ...props.teacher })

watch(() => props.teacher, (newVal) => {
  form.value = { ...newVal }
  previewImage.value = '' // reset preview ketika teacher berganti
})

const handleImageUpload = (event) => {
  const file = event.target.files[0]
  if (file) {
    const reader = new FileReader()
    reader.onload = (e) => {
      previewImage.value = e.target.result
    }
    reader.readAsDataURL(file)

    // Upload ke Supabase
    teachersStore.uploadTeacherImage(file)
      .then((response) => {
        form.value.profile_image = response.file_url || response.data?.file_url
      })
      .catch(() => {
        uiStore.showError('Failed to upload image')
      })
  }
}

const submitForm = async () => {
  try {
    await teacherService.updateTeacher(props.teacher.id, form.value)
    uiStore.showSuccess('Teacher updated successfully')
    emit('close')
  } catch (err) {
    console.error(err)
    uiStore.showError('Failed to update teacher')
  }
}
</script>

<style scoped>
.input {
  @apply w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500;
}
.btn-primary {
  @apply bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700;
}
.btn-secondary {
  @apply bg-gray-300 text-gray-800 px-4 py-2 rounded hover:bg-gray-400;
}
</style>
