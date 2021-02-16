<template>
      <v-card
    class="mx-auto"
    max-width="500"
    tile
  >
        <v-list dense>
          <v-subheader>ЕДИНИЦЫ ИЗМЕРЕНИЯ</v-subheader>
          <v-list-item-group
            color="primary"
          >
            <v-list-item
              v-for="(item, i) in items"
              :key="i"
              @click.stop="openDialog(item.id)">
       <!--       :to="'/unit/'+item.id"-->
              <v-list-item-content v-text="item.id" ></v-list-item-content>
              <v-list-item-content v-text="item.name" ></v-list-item-content>
              <v-list-item-content v-text="item.short_name"></v-list-item-content>
            </v-list-item>
          </v-list-item-group>
        </v-list>  
        <v-fab-transition>
          <v-btn
            color="pink"
            dark
            absolute
            bottom
            right
            fab
            @click.stop="openDialog(-1)"
          >
            <v-icon>mdi-plus</v-icon>
          </v-btn>
        </v-fab-transition>
        <Unit v-bind:content="mdata"/>
    </v-card> 
</template>
 
<script>
import Unit from "../components/Unit"
  export default {
    name: 'Units',
    components: {Unit},
    computed:{
      items(){
        return this.$store.getters.getUnits;
      }
    },
    data(){
      var mdata = new Object()
      mdata['dialog']=false
      mdata['id']=-1
      mdata['short_name']=""
      mdata['name']=""

      return {mdata:mdata}
    },
    created() {
      this.$store.dispatch('readUnits')
    },
    methods:{
      openDialog(id){
        this.mdata['id']=id
        if (id==-1) {
          this.mdata['short_name']=""
          this.mdata['name']=""
          this.mdata['dialog']=true        
        }else{
          this.$store.dispatch('readUnit', id)
            .then(resp=>{
              this.mdata['short_name']=resp.data.short_name
              this.mdata['name']=resp.data.name
              this.mdata['dialog']=true        
            })
            .catch(err => console.log(err))
        }
      },
    }
  }
</script>
