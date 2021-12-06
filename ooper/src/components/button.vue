<template>
  <div>
    <button type="submit" :disabled="disabled">{{ text }}</button>
  </div>
</template>

<script>
export default {
    props:["text","disabled"],
    data(){
        return{
            clientX:0,
            clientY:0
        }
    },
}
</script>

<style>
button{
    all:unset;
    cursor: pointer;
    border-radius: 50px;
    width:268px;
    height:68px;
    background:var(--bright-yellow);
    color:var(--purple);
    font-weight: bold;
    transition: all .1s linear;
}
button:hover{
  background:var(--dark-yellow);
}
button:disabled{
  cursor:default;
  background:var(--dark-yellow);
  
}
</style>