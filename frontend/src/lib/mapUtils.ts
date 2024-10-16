import { Icon, Style } from "ol/style";
import type { View } from "ol";
import type { Coordinate } from "ol/coordinate";

const musicalNoteSvg = '<svg width="120" height="120" xmlns="http://www.w3.org/2000/svg" class="ionicon" viewBox="0 0 512 512"><path d="M183.83 480a55.2 55.2 0 01-32.36-10.55A56.64 56.64 0 01128 423.58a50.26 50.26 0 0134.14-47.73L213 358.73a16.25 16.25 0 0011-15.49V92a32.1 32.1 0 0124.09-31.15l108.39-28.14A22 22 0 01384 54v57.75a32.09 32.09 0 01-24.2 31.19l-91.65 23.13A16.24 16.24 0 00256 181.91V424a48.22 48.22 0 01-32.78 45.81l-21.47 7.23a56 56 0 01-17.92 2.96z"/></svg>';
export const sessionStyle = new Style({
        image: new Icon({
            src: 'data:image/svg+xml;utf8,' + musicalNoteSvg,
            scale: 0.4,
        })
    })
    // new Style({
    //     image: new RegularShape({
    //         stroke: new Stroke(), //new Stroke({ color: [0, 0, 0, 0] }),
    //         points: 50,
    //         radius: 50, // <--------- control its size
    //         angle: Math.PI / 4
    //     })
    // })]

export const flyTo = (view: View, destination: Coordinate) => {
    const zoom = view.getZoom()!;
    const targetZoom = 15;
    view.animate(
        {
            center: destination,
        });
    if (zoom < targetZoom) {
        view.animate(
            {
                zoom: targetZoom,
            }
        );
    }
}