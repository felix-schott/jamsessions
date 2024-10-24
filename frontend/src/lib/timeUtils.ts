import type { Interval } from "../types";

// Helper func - given a Date object, returns the English common name for the day of the week.
let getDow = (d: Date) =>
    ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'][d.getDay()];

export const constructIntervalString = (interval: Interval, date: Date) => {
    let i: string;
    switch (interval) {
        case 'Once':
            i = 'as a one-off event';
            break;
        case 'Daily':
            i = 'everyday';
            break;
        case 'Weekly':
            i = `every week (${getDow(date)})`;
            break;
        case 'FirstOfMonth':
            i = 'every first ' + getDow(date) + ' of the month';
            break;
        case 'SecondOfMonth':
            i = 'every second ' + getDow(date) + ' of the month';
            break;
        case 'ThirdOfMonth':
            i = 'every third ' + getDow(date) + ' of the month';
            break;
        case 'FourthOfMonth':
            i = 'every fourth ' + getDow(date) + ' of the month';
            break;
        case 'LastOfMonth':
            i = 'every last ' + getDow(date) + ' of the month';
            break;
        default:
            throw "Unexpected value for property 'interval': " + interval;
    }
    return i
};

/** 
 * Calculate the number of minutes between two timestamps. 
 * It is a requirement for time2 to be after time1, and for the times to be provided in HH:MM (e.g. 15:30) format. 
 * */
export const minutesBetweenTimestamps = (time1: string, time2: string): number => {
    let d1 = new Date();
    const [h1, m1] = time1.split(":")
    console.log(parseInt(h1), parseInt(m1))
    d1.setHours(parseInt(h1), parseInt(m1), 0, 0)

    let d2 = new Date();
    const [h2, m2] = time2.split(":")
    d2.setHours(parseInt(h2), parseInt(m2), 0, 0)

    return Math.floor(
        (d2.getTime() - d1.getTime()) / 1000 / 60
    )
}