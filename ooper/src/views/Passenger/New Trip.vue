<template>
  <div id="map">
      <iframe
    width="680"
    height="600"
    frameborder="0" style="border:3px solid var(--bright-yellow);border-radius:50px;"
    :src="fullURL" allowfullscreen>
    </iframe>
  </div>
  <div id="form">
    <input v-model="origin" type="text" placeholder="start point"/>
    <input v-model="destination" type="text" placeholder="end point"/>
    <p>fee:<span id="value">{{ price }}</span></p>
    <Button text="book"/>
  </div>
</template>

<script>
import Button from "../../components/button.vue"
export default {
    components:{Button},
    data(){
        return{
            origin:"",
            destination:"",
            price:""
        }
    },
    computed:{
        fullURL(){
            if (this.origin != "" && this.destination != ""){
                return "https://www.google.com/maps/embed/v1/directions?key=" + process.env.VUE_APP_GMAPS_KEY + "&origin="+this.origin+"&destination="+this.destination
            }
            else{
                return "https://www.google.com/maps/embed/v1/view?key=" + process.env.VUE_APP_GMAPS_KEY +"&maptype=roadmap&center=1.38229,103.79714&zoom=10"
            }
        }
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