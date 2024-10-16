import { expect, test } from 'vitest'

import { minutesBetweenTimestamps } from './timeUtils'

test("minutesBetweenTimestamps returns 60 minutes", () => {
    expect(minutesBetweenTimestamps("15:00", "16:00")).toBe(60)
})