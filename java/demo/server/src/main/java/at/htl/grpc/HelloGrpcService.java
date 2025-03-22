package at.htl.grpc;

import io.quarkus.grpc.GrpcService;
import io.smallrye.mutiny.Multi;
import io.smallrye.mutiny.Uni;
import java.time.Duration;
import java.util.concurrent.atomic.AtomicReference;

@GrpcService
public class HelloGrpcService implements HelloWorldService {

    @Override
    public Uni<HelloReply> sayHello(HelloRequest request) {
        return Uni.createFrom()
            .item(
                HelloReply.newBuilder()
                    .setMessage(greet(request.getName()))
                    .build()
            );
    }

    public Multi<HelloReply> streamHello(Multi<HelloRequest> incomingStream) {
        AtomicReference<String> name = new AtomicReference<>("");
        incomingStream
            .subscribe()
            .with(request -> {
                name.set(request.getName());
            });
        return Multi.createFrom()
            .ticks()
            .every(Duration.ofSeconds(1))
            .map(ignored ->
                HelloReply.newBuilder().setMessage(greet(name.get())).build()
            );
    }

    private String greet(String name) {
        return String.format("Hello %s!", name);
    }
}
