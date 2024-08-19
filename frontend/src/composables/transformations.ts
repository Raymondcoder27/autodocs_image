import moment from "moment";

export function dateFormat(date:string) {
   if(date == null){
      return "--"
   }
   return moment(date).format("DD-MM-YYYY");
}

export function dateTimeFormat(date:string) {
   if(date == null){
      return "--"
   }
   return moment(date).format("DD-MM-YYYY HH:mm:ss");
}

export function timeLapse(date:string) {
   if(date == null){
      return "--"
   }

   let formattedDate = moment(date).format("DD-MM-YYYY HH:mm:ss");
   let startDate = moment(formattedDate, "DD-MM-YYYY HH:mm:ss")

   let elapsedDuration = moment.duration(moment().diff(startDate));

   let duration = 0
   let lapse = ""

   if(elapsedDuration.asSeconds() < 60){
      duration = Math.floor(elapsedDuration.asSeconds())
      lapse = (duration >  1) ? " secs ago" :" sec ago"

      return duration + lapse
   }
   if(elapsedDuration.asMinutes() < 60){
      duration = Math.floor(elapsedDuration.asMinutes())
      lapse = (duration >  1) ? " mins ago" : " min ago"

      return duration + lapse
   }

   if(elapsedDuration.asHours() < 24){
      duration = Math.floor(elapsedDuration.asHours())
      lapse = (duration >  1) ? " hrs ago" : " hr ago"

      return duration + lapse
   }

   if(elapsedDuration.asDays() > 1 && elapsedDuration.asDays() < 7){
      duration = Math.floor(elapsedDuration.asDays())
      lapse =  (duration >  1) ? " days ago" : " day ago"

      return duration + lapse
   }

   if(elapsedDuration.asWeeks() > 0 && elapsedDuration.asWeeks() < 4){
      duration = Math.floor(elapsedDuration.asWeeks())
      lapse =  (duration >  1) ? " weeks ago" : " week ago"

      return duration + lapse
   }

   return formattedDate


}

export function percentage(numerator:number|undefined, denominator:number|undefined) {
   if(numerator == null || denominator == null){
      return 0;
   }
   return ((numerator/denominator) * 100).toFixed(2);
}