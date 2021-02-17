<template>
    <v-dialog
      v-model="this.content.editPrice"
      persistent
      max-width="500"
    >
       <v-card>
        <v-card-title>
          <span class="headline">Цена</span>
        </v-card-title>
        <v-card-text>
          <v-text-field
            label="Наименование"
            v-model="content.name"
            readonly
          ></v-text-field>  
          <v-container>
            <v-row>
              <v-col
                cols="12"
                sm="6"
                md="4"
              >
                <v-text-field
                  label="Ед. изм."
                  readonly
                  v-model="content.price_unit_short_name"
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
                  v-model.number="content.price"
                  :rules="[rules.num]"
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
            @click="content.editPrice = false"
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
    name: 'MaterialPrice',
    props:{content:Object
    },
    data(){
      return {rules:{num: value => {return !isNaN(value)||'Должно быть число'}}}
    },
    methods:{
      saveData(){
        if (isNaN(this.content.price)){
          return
        }
        this.$store.dispatch('writeMaterialPrice',{id:-1,material_id:this.content.id,price:this.content.price})
        this.content.editPrice= false

      }
    }
  }
</script>
