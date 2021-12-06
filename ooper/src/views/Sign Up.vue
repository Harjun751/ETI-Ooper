<template>
  <h1>new to ooper?</h1>
    <div id="toggle"><Toggle @updatePassenger="updatePassenger"/></div>
  <div class="row">
      <input type="text" v-model="firstName" placeholder="first name" required/>
      <input type="text" v-model="lastName" placeholder="last name" required/>
  </div>
  <div class="row">
      <input type="text" v-model="mobileNumber" placeholder="mobile number" required/>
      <input type="email" v-model="email" placeholder="e-mail address" required/>
  </div>
  <div class="row" id="driver-only" :class="{hide:isPassenger,show:isPassenger==false}">
      <input type="text" v-model="ic" placeholder="identification number" required :disabled="isPassenger"/>
      <input type="email" v-model="license" placeholder="car license number" required :disabled="isPassenger"/>
  </div>
  <div class="row">
      <input type="password" v-model="password" placeholder="password" required/>
  </div>
  <Button text="sign up" @click="submitSignUp"/>
</template>

<script>
import Toggle from "../components/toggle.vue"
import Button from "../components/button.vue"
const Swal = require('sweetalert2')
export default {
    components:{Toggle,Button},
    data(){
        return{
            isPassenger:true,
            firstName:"",
            lastName:"",
            mobileNumber:null,
            email:"",
            ic:"",
            license:"",
            password:""
        }
    },
    methods:{
        updatePassenger(passenger){
            this.isPassenger = passenger
        },
        async submitSignUp(){
            var data = {"FirstName":this.firstName,"LastName":this.lastName,"MobileNumber":Number(this.mobileNumber),"Email":this.email,"Password":this.password}
            var url = ""
            if (this.isPassenger){
                url = process.env.VUE_APP_PASSENGER_MS_HOST+"/api/v1/passengers"
            }
            else{
                url = process.env.VUE_APP_DRIVER_MS_HOST+"/api/v1/drivers"
                data["ICNumber"] = this.ic
                data["LicenseNumber"] = this.license
            }
            await fetch(url,{
                body: JSON.stringify(data),
                method:"POST",
                headers: {
                    'Content-Type': 'application/json',
                },
            }).then((resp)=>{
                if (resp.status==200) {
                    Swal.fire({
                        title: 'done!',
                        text: 'account created',
                        icon: 'success',
                        confirmButtonText: 'close',
                        customClass:{
                            popup: 'custom-swal-modal',
                            icon: 'custom-swal-icon',
                            content: 'custom-swal-content',
                            confirmButton: 'custom-swal-button'
                        }
                    })
                }
            })
        }
    },
}
</script>

<style scoped>
h1{
    font-size:54px;
    color:var(--dark-yellow);
    font-weight:normal;
}
.row{
    display:flex;
    justify-content: center;
}
input{
    margin:0 30px 30px 30px;
}
#toggle{
    display: flex;
    justify-content: center;
    margin-top:100px;
}
::v-deep(button){
    margin-top:40px;
}
#driver-only{
    transition: height 0.1s;
}
#driver-only.show{
    height:74px;
}
#driver-only.hide{
    height:0px;
}
</style>