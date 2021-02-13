<template>
      <v-card
    class="mx-auto"
    max-width="300"
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
              :to="'/unit/'+item.id"
            >
              <v-list-item-content>
                <v-list-item-title v-text="item.id"></v-list-item-title>
                <v-list-item-title v-text="item.name"></v-list-item-title>
                <v-list-item-title v-text="item.short_name"></v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list-item-group>
        </v-list>    
    </v-card> 
</template>
 
<script>
  export default {
    name: 'Units',
    computed:{
      items(){
        return this.$store.state.units
      }
    }
  }
</script>
