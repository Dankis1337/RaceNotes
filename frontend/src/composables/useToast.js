import { ref } from 'vue'

const message = ref('')
const type = ref('success')
const show = ref(false)

export function useToast() {
  function toast(msg, t = 'success') {
    message.value = msg
    type.value = t
    show.value = true
  }

  function close() {
    show.value = false
  }

  return { message, type, show, toast, close }
}
