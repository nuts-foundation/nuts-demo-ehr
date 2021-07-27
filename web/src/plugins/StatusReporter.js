import {ref} from 'vue'

export default {
    install: (app) => {
        let errorMessage = ref("")
        let statusMessage = ref("")

        app.config.globalProperties.$status = {
            error: (newMessage) => {
                console.error('An error occurred: ' + newMessage)
                errorMessage.value = newMessage
            },
            clearError: () => {
                errorMessage.value = ""
            },

            status: (newMessage) => {
               statusMessage.value = newMessage
            },
            clearStatus: () => {
                statusMessage.value = ""
            },
            errorMessage: errorMessage,
            statusMessage: statusMessage
        }
    }
}
