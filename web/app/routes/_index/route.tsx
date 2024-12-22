import type { MetaFunction } from "@remix-run/node";

export const meta: MetaFunction = () => {
    return [
        { title: "Hastypal" },
        { name: "description", content: "Welcome to Hastypal!" },
    ];
};

export default function Index() {
    return (<div>Hastypal!</div>);
}
