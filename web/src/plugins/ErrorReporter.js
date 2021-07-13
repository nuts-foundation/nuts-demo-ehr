import {ref} from 'vue'

export default {
    install: (app) => {
        let message = ref("")
        app.config.globalProperties.$errors = {
            report: (newMessage) => {
                console.error('An error occurred: ' + newMessage)
                message.value = newMessage
            },
            clear: () => {
                message.value = ""
            },
            message: message,
        }
    }
}
