export const lessonTimeStart = {
  "09:00": 1,
  "10:45": 2,
  "13:00": 3,
  "14:45": 4,
  "16:30": 5,
  "18:15": 6,
  "20:00": 7,
};

export const lessonTimeEnd = {
  1: "10:30",
  2: "12:15",
  3: "14:30",
  4: "16:15",
  5: "18:00",
  6: "19:45",
  7: "21:30",
};

export function getLessonNumber(lessons) {
  let lessonsWithLessonNumber = {};
  if (!lessons) return 0;
  lessons.forEach(function (lesson) {
    let lessonNumber = lessonTimeStart[lesson.time_start];
    lessonsWithLessonNumber[lessonNumber] = lesson;
  });

  return lessonsWithLessonNumber;
}
