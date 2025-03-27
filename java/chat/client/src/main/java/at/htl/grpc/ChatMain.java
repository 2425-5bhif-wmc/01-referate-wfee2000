package at.htl.grpc;

import at.htl.grpc.ChatOuterClass.ClaimNameRequest;
import at.htl.grpc.ChatOuterClass.IncomingMessage;
import at.htl.grpc.ChatOuterClass.OutgoingMessage;
import io.grpc.Metadata;
import io.grpc.Metadata.Key;
import io.quarkus.grpc.GrpcClient;
import io.quarkus.grpc.GrpcClientUtils;
import io.quarkus.runtime.QuarkusApplication;
import io.quarkus.runtime.annotations.QuarkusMain;
import io.smallrye.mutiny.Multi;
import io.smallrye.mutiny.subscription.MultiEmitter;
import java.io.IOException;
import org.jline.terminal.Terminal;
import org.jline.terminal.TerminalBuilder;

@QuarkusMain
public class ChatMain implements QuarkusApplication {

    @GrpcClient
    Chat chat;

    private static StringBuilder message = new StringBuilder();

    @Override
    public int run(String... args) throws Exception {
        String token = chat
            .claimName(
                ClaimNameRequest.newBuilder()
                    .setName(String.join(" ", args))
                    .build()
            )
            .await()
            .indefinitely()
            .getToken();

        Metadata headers = new Metadata();
        headers.put(
            Key.of("authorization", Metadata.ASCII_STRING_MARSHALLER),
            String.format("Bearer %s", token)
        );

        Chat authorizedStub = GrpcClientUtils.attachHeaders(chat, headers);

        Multi<OutgoingMessage> outgoingStream = Multi.createFrom()
            .emitter(em -> {
                SendMessages(em);
            });

        PrintMessages(authorizedStub.connect(outgoingStream));

        return 0;
    }

    public static void PrintMessages(Multi<IncomingMessage> incomingStream) {
        incomingStream
            .subscribe()
            .with(incomingMessage -> {
                System.out.printf(
                    "\033[2K\r%s: %s\n\rWrite message: %s",
                    incomingMessage.getName(),
                    incomingMessage.getResponse(),
                    message
                );
            });
    }

    public static void SendMessages(
        MultiEmitter<? super OutgoingMessage> outgoingStream
    ) {
        //disableCanonical();

        while (true) {
            try {
                readMessage();
            } catch (IOException e) {
                throw new RuntimeException(e);
            }

            outgoingStream.emit(
                OutgoingMessage.newBuilder().setMessage(message.toString()).build()
            );

            System.out.printf("\033[A\033[2K\rYou wrote: %s\n\r", message);

            message.setLength(0);
        }
    }

    private static void readMessage() throws IOException {
        System.out.print("Write message: ");

        try (Terminal terminal = TerminalBuilder.terminal()) {
            terminal.enterRawMode();
            int controlCharCounter = 0;

            while (true) {
                char character = (char) terminal.reader().read();

                if (controlCharCounter == 2 && character == '[') {
                    controlCharCounter--;
                    continue;
                }

                if (controlCharCounter == 1) {
                    controlCharCounter--;
                    continue;
                }


                if (character == '\n' || character == '\r') {
                    System.out.print("\n");
                    break;
                }

                if (character == 127) {
                    if (message.length() > 1) {
                        System.out.print("\b \b");
                        message.deleteCharAt(message.length() - 1);
                    } else if (message.length() == 1) {
                        System.out.print("\r\033[2K\rWrite message: ");
                    }
                }

                controlCharCounter = 0;

                if (Character.isISOControl(character)) {
                    controlCharCounter = 2;
                    continue;
                }

                System.out.print(character);
                message.append(character);
            }
        }
    }
}
