package at.htl.grpc.chat;

import io.grpc.Grpc;
import io.grpc.InsecureChannelCredentials;
import io.grpc.ManagedChannel;
import io.grpc.stub.StreamObserver;

public class ChatStub {
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
    }
}
