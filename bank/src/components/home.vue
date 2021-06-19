<template>
  <div>
      <!-- <p>{{transactions}}</p> -->
<table border=1>
  <thead>
    <tr>
      <th>ID</th>
      <th>Date</th>
      <th>Description</th>
      <th>Ref</th>
      <th>Debit</th>
      <th>Credit</th>
      <th>Balance</th>
    </tr>
  </thead>
  <tbody v-for="item in transactions" :key="item.id">
    <tr>
      <td>{{ item.id }}</td>
      <td>{{ item.date | formatDate }}</td>
      <td>{{ item.description }}</td>
      <td>{{ item.ref.String }}</td>
      <td>{{ item.debit.Float64 }}</td>
      <td>{{ item.credit.Float64 }}</td>
      <td>{{ item.bal }}</td>
    </tr>
  </tbody>
</table>
  </div>
</template>
<script>
export default {
  name: "Index",
  data() {
    return {
      transactions: null
    };
  },
  mounted() {
    this.$http
      .get("/api/all")
      .then((response) => {
        this.transactions = response.data;
      })
      .catch((error) => console.log(error));
  },
};
</script>
<style>
table {
  width: 100%;
  background-color: #ffffff;
  border-collapse: collapse;
  border-width: 2px;
  border-color: #ffcc00;
  border-style: solid;
  color: #000000;
}
</style>