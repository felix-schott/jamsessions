import { expect, test } from 'vitest'

import { sanitisePathElement, extractDomain, processVenueAndSessionUrl } from './uriUtils';

// sanitisePath

test("sanitisePath removes special characters", () => {
    expect(sanitisePathElement("foo'sbar&restaurant")).toBe("foosbarrestaurant");
})

test("sanitisePath outputs lower case", () => {
    expect(sanitisePathElement("Foosbar")).toBe("foosbar");
})

test("sanitisePath outputs replaces white space with hyphen", () => {
    expect(sanitisePathElement("foos bar restaurant")).toBe("foos-bar-restaurant");
})

test("sanitisePath formats correctly", () => {
    expect(sanitisePathElement("Foo's Bar")).toBe("foos-bar");
})

// extractDomain

test("extractDomain works correctly (http)", () => {
    expect(extractDomain("http://business.instagram.com/?locale=en_GB/")).toBe("business.instagram.com");
})

test("extractDomain works correctly (https)", () => {
    expect(extractDomain("https://business.instagram.com/?locale=en_GB/")).toBe("business.instagram.com");
})

test("extractDomain works correctly (no path)", () => {
    expect(extractDomain("https://business.instagram.com/")).toBe("business.instagram.com");
})

test("extractDomain works correctly (no path, no trailing slash)", () => {
    expect(extractDomain("https://business.instagram.com")).toBe("business.instagram.com");
})

// process venue and session url

test("processVenueAndSessionUrl works correctly with same domain, different paths", () => {
    expect(processVenueAndSessionUrl("http://example.com", "http://example.com/events/1")).toStrictEqual({
        venueWebsite: "example.com",
        sessionWebsite: "example.com/events"
    });
})

test("processVenueAndSessionUrl works correctly with same domain", () => {
    expect(processVenueAndSessionUrl("http://example.com", "http://example.com")).toStrictEqual({
        venueWebsite: "example.com",
        sessionWebsite: null
    });
})

test("processVenueAndSessionUrl works correctly with different domains", () => {
    expect(processVenueAndSessionUrl("http://example.com", "http://example2.com")).toStrictEqual({
        venueWebsite: "example.com",
        sessionWebsite: "example2.com"
    });
})