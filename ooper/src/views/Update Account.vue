<template>
  <h1>...oops?</h1>
  <div class="row">
      <input type="text" v-model="firstName" placeholder="first name" required/>
      <input type="text" v-model="lastName" placeholder="last name" required/>
  </div>
  <div class="row">
      <input type="text"  v-model="mobileNumber" placeholder="mobile number" required/>
      <input type="email" v-model="email" placeholder="e-mail address" required/>
  </div>
  <div class="row" v-if="isPassenger==false">
      <input type="email" v-model="license" placeholder="car license number" required/>
  </div>
  <Button text="update" @click="updateUser" />
</template>

<script>
import Button from "../components/button.vue"
import { store } from "../state"
const Swal = require('sweetalert2')
export default {
    components:{Button},
    data(){
        return{
            isPassenger:store.state.isPassenger,
            firstName:"",
            lastName:"",
            mobileNumber:null,
            license:"",
            email:""
        }
    },
    methods:{
        async updateUser(){
            var data = {"FirstName":this.firstName,"LastName":this.lastName,"MobileNumber":Number(this.mobileNumber),"Email":this.email}
            var url = ""
            if (this.isPassenger){
                url = process.env.VUE_APP_PASSENGER_MS_HOST+"/api/v1/passengers"
            }
            else if (this.isPassenger == false){
                url = process.env.VUE_APP_DRIVER_MS_HOST+"/api/v1/drivers"
                data["LicenseNumber"] = this.license
            }
            await fetch(url,{
                body: JSON.stringify(data),
                method:"PATCH",
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': "Bearer " +  store.state.jwtAccessToken
                },
            }).then((resp)=>{
                if (resp.status==200) {
                    Swal.fire({
                        title: 'done!',
                        text: 'your account has been updated',
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
    }
}
</script>

<style scoped>
h1{
    font-size:54px;
    color:var(--dark-yellow);
    font-weight:normal;
    margin-top:150px;
}
.row{
    display:flex;
    justify-content: center;
}
input{
    margin:0 30px 30px 30px;
}
::v-deep(button){
    margin-top:40px;
}
</style>