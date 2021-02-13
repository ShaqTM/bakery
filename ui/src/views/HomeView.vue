<template>
      <v-card
    class="mx-auto"
    max-width="300"
    tile
  >
        <v-list dense>
          <v-subheader>MENU</v-subheader>
          <v-list-item-group
            v-model="selectedItem"
            color="primary"
          >
            <v-list-item
              v-for="(item, i) in items"
              :key="i"
              :to="item.link"
            >
              <v-list-item-content>
                <v-list-item-title v-text="item.text"></v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list-item-group>
        </v-list>    
    </v-card> 
</template>

<script>
  export default {
    name: 'HomeView',
    data: () => ({
      selectedItem: 1,
      items: [
        { text: 'Заказы',link:'/orders'},
        { text: 'Рецепты',link:'/recipes'},
        { text: 'Материалы',link:'/materials'},
        { text: 'Единицы измерения',link:'/units'},
      ],
    }),
  }
</script>
