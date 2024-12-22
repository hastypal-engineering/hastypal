import type { MetaFunction } from "@remix-run/node";

export const meta: MetaFunction = () => {
    return [
        { title: "Hastypal" },
        { name: "description", content: "Welcome to Hastypal!" },
    ];
};

export default function Index() {
    return (
        <div>
            <div>
                <p>Hastypal</p>
                <button>Entrar</button>
                <button>Registrarse</button>
            </div>
            <div>
                <h1>Hastypal ayuda a tus clientes a contratar tus servicios</h1>
                <h3>Una forma más fácil de vender tus servicios con ayuda de nuevas herramientas</h3>
                <button>Empieza ahora</button>
            </div>
            <div>
                <h2>Precios</h2>
                <div>
                    <h4>100€/año</h4>
                </div>
            </div>
        </div>
    );
}
