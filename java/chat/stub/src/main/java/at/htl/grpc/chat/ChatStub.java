package at.htl.grpc.chat;

import io.grpc.Grpc;
import io.grpc.InsecureChannelCredentials;
import io.grpc.ManagedChannel;
import io.grpc.stub.StreamObserver;

import java.io.IOException;
import java.io.InputStream;

public class ChatStub {
    private static String message = "";

    public static void main(String[] args) {
        ManagedChannel channel = Grpc
                .newChannelBuilder("localhost:5555", InsecureChannelCredentials.create())
                .build();
        ChatServiceGrpc.ChatServiceBlockingStub messageStub = ChatServiceGrpc.newBlockingStub(channel);
        String token = messageStub.claimName(Chat.ClaimNameRequest.newBuilder().setName("Winnie").build()).getToken();

        ChatServiceGrpc.ChatServiceStub chatStub = ChatServiceGrpc.newStub(channel);
        StreamObserver<Chat.OutgoingMessage> outgoingObserver = chatStub
                .withCallCredentials(new AuthorizationCallCredentials(token))
                .connect(new StreamObserver<>() {
                    @Override
                    public void onNext(Chat.IncomingMessage incomingMessage) {
                        System.out.printf(
                                "\033[2K\r%s: %s\n\rWrite message: %s",
                                incomingMessage.getName(),
                                incomingMessage.getResponse(),
                                message);
                    }

                    @Override
                    public void onError(Throwable throwable) {
                        throwable.printStackTrace();
                    }

                    @Override
                    public void onCompleted() {
                        // should never happen
                    }
                });

        try {
            SendMessages(outgoingObserver);
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }

    public static void SendMessages(StreamObserver<Chat.OutgoingMessage> outgoingObserver) throws IOException {
        disableCanonical(System.in);

        while (true) {
            System.out.print("Write message: ");
            while (true) {
                char character = (char) System.in.read();

                if (character == '\n' || character == '\r') {
                    break;
                }

                message = message + character;
                System.out.print("In");
            }

            outgoingObserver.onNext(Chat.OutgoingMessage.newBuilder().setMessage(message).build());
            System.out.printf("\033[A\033[2K\rYou wrote: %s\n\r", message);

            message = "";
        }
    }

    private static void disableCanonical(InputStream inputStream) {
        // TODO: disable Canonical for inputStream
    }
}
