<template>
<div>
  <v-tabs
    v-model="tab"
    fixed-tabs
    background-color="indigo"
    dark
    @change="tabChanged"
    >
    <v-tab 
      v-for = "item in items"
      :key = "item.text">
      {{ item.text }}      
    
    </v-tab>
  </v-tabs>
  <v-tabs-items v-model="tab">
    <v-tab-item
      v-for="item in items"
      :key="item.text"
    >
    <v-card flat>
      <component v-bind:is="item.component"></component>
          
        </v-card>
      </v-tab-item>
    </v-tabs-items>  
  </div>
</template>
<script>
  import Orders from "./Orders.vue"
  import Units from "./Units.vue"
  import Recipes from "./Recipes.vue"
  import Materials from "./Materials.vue"

  export default {
    name: 'HomeView',
    components: {Materials,Recipes,Orders,Units},
    data: () => ({
      tab: 'Заказы',
      items: [
        { text: 'Единицы измерения',component:Units},
        { text: 'Материалы',component:Materials},
        { text: 'Рецепты',component:Recipes},
        { text: 'Заказы',component:Orders}
        ],      
    }),
    methods:{
      tabChanged(){
        if (this.tab==3){
          this.$store.dispatch('readOrders')
        }else if (this.tab==1){
          this.$store.dispatch('readMaterials')
        }else if (this.tab==0){
          this.$store.dispatch('readUnits')
        }else if (this.tab==2){
          this.$store.dispatch('readRecipes')
        }



      }
    }
  }
</script>
