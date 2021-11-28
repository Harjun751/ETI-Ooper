<template>
    <div v-if="exists">
        <div id="map">
            <iframe
            width="680"
            height="600"
            frameborder="0" style="border:3px solid var(--bright-yellow);border-radius:50px;"
            :src="fullURL" allowfullscreen>
            </iframe>
        </div>
        <div id="form">
            <input v-model="origin" type="text" placeholder="start point" disabled/>
            <input v-model="destination" type="text" placeholder="end point" disabled/>
            <p>fee:<span id="value">{{ price }}</span></p>
            <Button text="start trip" @click="startTrip"/>
            <Button text="end trip" @click="endTrip"/>
        </div>
    </div>
    <div v-if="exists==false">
        <h3 style="color:var(--bright-yellow);">No Trips!</h3>
    </div>
</template>

<script>
import Button from "../../components/button.vue"
import { store } from "../../state"
export default {
    components:{Button},
    data(){
        return{
            origin:"",
            destination:"",
            price:"",
            start:"",
            end:"",
            tripID:"",
            exists:null
        }
    },
    computed:{
        fullURL(){
            return "https://www.google.com/maps/embed/v1/directions?key=" + process.env.VUE_APP_GMAPS_KEY + "&origin="+this.origin+"&destination="+this.destination
        }
    },
    async mounted(){
        await fetch("http://localhost:5004/api/v1/current-trip",{
        method:"GET",
        headers: {
            'Content-Type': 'application/json',
            'Authorization': "Bearer " +  store.state.jwtAccessToken
        },
        })
        .then(async (res)=>{
            if (res.status==404){
                this.exists = false
                throw "No data"
            }
            await res.json()
        })
        .then((data)=>{
            this.exists = true
            this.origin = data.PickUp
            this.destination = data.DropOff
            this.start = data.Start
            this.end = data.End
            this.tripID = data.ID
        })
    },
    methods:{
        async startTrip(){
            await fetch("http://localhost:5004/api/v1/trip/"+this.tripID+"/start",{
                method:"POST",
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': "Bearer " +  store.state.jwtAccessToken
                },
            })
        },
        async endTrip(){
            await fetch("http://localhost:5004/api/v1/trip/"+this.tripID+"/end",{
                method:"POST",
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': "Bearer " +  store.state.jwtAccessToken
                },
            })
        },
    }
}
</script>

<style scoped>
#form{
    float:right;
    margin-top:200px;
    margin-right:30px;
}
#form p{
    color:var(--bright-yellow);
    text-align: left;
    margin-bottom:50px;
}
#value{
    padding-left:10px;
    text-decoration: underline;
}
#map{
    float:left;
    margin-top:100px;
    margin-left:150px;
}
input{
    width:700px;
}
</style>