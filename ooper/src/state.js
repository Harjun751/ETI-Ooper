import { reactive } from "@vue/reactivity"
// eslint-disable-next-line no-unused-vars
const store = {
    state: reactive({
      loggedIn:false,
      isPassenger:null
    }),
  
    setLoggedIn(newValue) {
        this.state.loggedIn = newValue
    },
  
    setIsPassenger(newValue){
        this.state.isPassenger = newValue
    }
  }
export {store};