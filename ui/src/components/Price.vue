<template>
    <v-dialog
      v-model="this.content.editPrice"
      max-width="500"
    >
       <v-card>
        <v-card-title>
          <span class="headline">Цена</span>
        </v-card-title>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col
                cols="12"
                sm="6"
                md="4"
              >
                <v-text-field
                  label="Наименование"
                  required
                  v-model="content.name"
                ></v-text-field>
              </v-col>
              <v-col
                cols="12"
                sm="6"
                md="4"
              >
                <v-text-field
                  label="Цена"
                  required
                  v-model="content.price"
                ></v-text-field>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="blue darken-1"
            text
            @click="content.dialog = false"
          >
            Закрыть
          </v-btn>
          <v-btn
            color="blue darken-1"
            text
            @click="saveData()"
          >
            Сохранить
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
</template>
 
<script>
  export default {
    name: 'Price',
    props:{content:Object
    },

    methods:{
      saveData(){
        this.$store.dispatch('writePrice',{id:-1,material_id:this.content.id,price:this.content.price})
        this.content.editPrice= false

      }
    }
  }
</script>
