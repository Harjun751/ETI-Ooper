import { reactive } from "@vue/reactivity"
// eslint-disable-next-line no-unused-vars
const store = {
    state: reactive({
      jwtAccessToken:null,
      isPassenger:null
    }),
  
    setJWTAccessToken(newValue) {
        this.state.jwtAccessToken = newValue
    },
  
    setIsPassenger(newValue){
        this.state.isPassenger = newValue
    },

  }
export {store};