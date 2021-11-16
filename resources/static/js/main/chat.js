function formatDate(date) {

    let dd = date.getDate();
    if (dd < 10) dd = '0' + dd;

    let mm = date.getMonth() + 1;
    if (mm < 10) mm = '0' + mm;

    let yy = date.getFullYear() % 100;
    if (yy < 10) yy = '0' + yy;

    let hh = date.getHours();
    if (hh < 10) hh = '0' + hh;

    let mi = date.getMinutes();
    if (mi < 10) mi = '0' + mi;

    return dd + '.' + mm + '.' + yy + ' - ' + hh + ':' + mi;
}

Vue.component('vue-date', {
    props: {
        v: Number,
    },
    data: function() {
        let date = new Date(this.v * 1000);
        return {
            r: formatDate(date),
        }
    },
    template: '<div class="time_sent">{{ r }}</div>'
});

new Vue({
    el: "#app",
    data: {
        show: false
    }
})