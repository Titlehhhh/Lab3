<template>
  <section class="container my-5">
    <div class="text-center">
      <h1>Добро пожаловать в химчистку “Чисто и Точка”</h1>
      <p class="lead">Профессиональная чистка одежды, обуви и текстиля в вашем городе.</p>
      <RouterLink to="/products" class="btn btn-primary">Наши услуги</RouterLink>
    </div>
  </section>

  <section class="container my-5">
    <div class="row text-center">
      <div class="col-md-4">
        <h3>Одежда</h3>
        <p>Удалим любые пятна и вернём свежесть вашей одежде.</p>
      </div>
      <div class="col-md-4">
        <h3>Обувь</h3>
        <p>Качественная чистка и восстановление обуви.</p>
      </div>
      <div class="col-md-4">
        <h3>Текстиль</h3>
        <p>Чистим шторы, ковры и постельное бельё.</p>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>
