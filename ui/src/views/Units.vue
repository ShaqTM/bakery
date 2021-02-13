<template>
      <v-card
    class="mx-auto"
    max-width="300"
    tile
  >
        <v-list dense>
          <v-subheader>ЕДИНИЦЫ ИЗМЕРЕНИЯ</v-subheader>
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
    name: 'Units',
    data: () => ({
      selectedItem: 1,
      items: [
        { text: 'Единица 1',link:'/unit/1'},
        { text: 'Единица 2',link:'/unit/2'},
        { text: 'Единица 3',link:'/unit/3'},
        { text: 'Единица 4',link:'/unit/4'},
      ],
    }),
  }
</script>
